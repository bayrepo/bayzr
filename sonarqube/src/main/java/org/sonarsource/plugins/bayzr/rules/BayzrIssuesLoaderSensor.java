package org.sonarsource.plugins.bayzr.rules;

import java.io.File;
import java.util.Arrays;
import java.util.List;
import java.sql.*;

import javax.xml.stream.XMLStreamException;

import org.apache.commons.lang.StringUtils;
import org.sonar.api.batch.fs.FileSystem;
import org.sonar.api.batch.fs.InputFile;
import org.sonar.api.batch.sensor.Sensor;
import org.sonar.api.batch.sensor.SensorContext;
import org.sonar.api.batch.sensor.SensorDescriptor;
import org.sonar.api.batch.sensor.issue.NewIssue;
import org.sonar.api.batch.sensor.issue.NewIssueLocation;
import org.sonar.api.config.Settings;
import org.sonar.api.rule.RuleKey;
import org.sonar.api.utils.log.Logger;
import org.sonar.api.utils.log.Loggers;
import org.sonarsource.plugins.bayzr.languages.CppLanguage;

/**
 * The goal of this Sensor is to load the results of an analysis performed by a fictive external tool named: FooLint
 * Results are provided as an xml file and are corresponding to the rules defined in 'rules.xml'.
 * To be very abstract, these rules are applied on source files made with the fictive language Foo.
 */
public class BayzrIssuesLoaderSensor implements Sensor {

  private static final Logger LOGGER = Loggers.get(BayzrIssuesLoaderSensor.class);

  protected static final String REPORT_USER_KEY = "sonar.bayzr.user";
  protected static final String REPORT_PASS_KEY = "sonar.bayzr.pass";
  protected static final String REPORT_URL_KEY = "sonar.bayzr.url";

  protected final Settings settings;
  protected final FileSystem fileSystem;
  protected SensorContext context;

  /**
   * Use of IoC to get Settings, FileSystem, RuleFinder and ResourcePerspectives
   */
  public BayzrIssuesLoaderSensor(final Settings settings, final FileSystem fileSystem) {
    this.settings = settings;
    this.fileSystem = fileSystem;
  }

  @Override
  public void describe(final SensorDescriptor descriptor) {
    descriptor.name("BayZR Issues Loader Sensor");
    descriptor.onlyOnLanguage(CppLanguage.KEY);
  }

  protected String reportUrlKey() {
    return REPORT_URL_KEY;
  }

  protected String getReportURL() {
    String reportUrl = settings.getString(reportUrlKey());
    if (!StringUtils.isEmpty(reportUrl)) {
      return reportUrl;
    } else {
      return "jdbc:mysql://localhost/bayzr";
    }
  }

  protected String reportUserKey() {
    return REPORT_USER_KEY;
  }

  protected String getReportUser() {
    String reportUser = settings.getString(reportUserKey());
    if (!StringUtils.isEmpty(reportUser)) {
      return reportUser;
    } else {
      return "bayzr";
    }
  }

  protected String reportPassKey() {
    return REPORT_PASS_KEY;
  }

  protected String getReportPass() {
    String reportPass = settings.getString(reportPassKey());
    if (!StringUtils.isEmpty(reportPass)) {
      return reportPass;
    } else {
      return "bayzr";
    }
  }

  @Override
  public void execute(final SensorContext context) {
      if (!StringUtils.isEmpty(getReportUrl())&&!StringUtils.isEmpty(getReportUser())&&!StringUtils.isEmpty(getReportPass())){
       this.context = context;
       try {
        String url = getReportUrl();
        String user = getReportUser();
        String pass = getReportPass();
        parseAndSaveResults(url, user, pass);
       } catch (XMLStreamException e) {
        throw new IllegalStateException("Unable to parse the provided BayZR info", e);
       }
      } else {
        LOGGER.info("Empty URL or User or Password parameter");
      }
  }

  protected void parseAndSaveResults(String url, String user, String pass) throws XMLStreamException {
    LOGGER.info("Parsing 'BayZR' Analysis Results");
    BayzrAnalysisResultsParser parser = new BayzrAnalysisResultsParser();
    List<BayzrError> errors = parser.parse(url, user, pass);
    for (BayzrError error : errors) {
      getResourceAndSaveIssue(error);
    }
  }

  private void getResourceAndSaveIssue(final BayzrError error) {
    LOGGER.debug(error.toString());

    InputFile inputFile = fileSystem.inputFile(
      fileSystem.predicates().and(
        fileSystem.predicates().hasRelativePath(error.getFilePath()),
        fileSystem.predicates().hasType(InputFile.Type.MAIN)));

    LOGGER.info("inputFile null ? " + inputFile);

    if (inputFile != null) {
      saveIssue(inputFile, error.getLine(), error.getType(), error.getDescription());
    } else {
      LOGGER.error("Not able to find a InputFile with " + error.getFilePath());
    }
  }

  private void saveIssue(final InputFile inputFile, int line, final String externalRuleKey, final String message) {
    RuleKey ruleKey = RuleKey.of(BayzrRulesDefinition.getRepositoryKeyForLanguage(inputFile.language()), externalRuleKey);

    NewIssue newIssue = context.newIssue()
      .forRule(ruleKey);

    NewIssueLocation primaryLocation = newIssue.newLocation()
      .on(inputFile)
      .message(message);
    if (line > 0) {
      primaryLocation.at(inputFile.selectLine(line));
    }
    newIssue.at(primaryLocation);

    LOGGER.info("Issue " + newIssue);
    newIssue.save();
  }

  @Override
  public String toString() {
    return "BayzrIssuesLoaderSensor";
  }

  private class BayzrError {

    private final String type;
    private final String description;
    private final String filePath;
    private final int line;

    public BayzrError(final String type, final String description, final String filePath, final int line) {
      this.type = type;
      this.description = description;
      this.filePath = filePath;
      this.line = line;
    }

    public String getType() {
      return type;
    }

    public String getDescription() {
      return description;
    }

    public String getFilePath() {
      return filePath;
    }

    public int getLine() {
      return line;
    }

    @Override
    public String toString() {
      StringBuilder s = new StringBuilder();
      s.append(type);
      s.append("|");
      s.append(description);
      s.append("|");
      s.append(filePath);
      s.append("(");
      s.append(line);
      s.append(")");
      return s.toString();
    }
  }

  private class BayzrAnalysisResultsParser {

    public List<BayzrError> parse(String url, String user, String pass) throws XMLStreamException {
      List<BayzrError> fndIssuesList = new ArrayList<BayzrError>();
      Connection conn = null;
      Statement stmt = null;
      LOGGER.info("Parsing file {}");
      try{
        Class.forName("com.mysql.jdbc.Driver");
        LOGGER.info("Connecting to database...");
        conn = DriverManager.getConnection(url,user,pass);

        LOGGER.info("Creating statement...");
        stmt = conn.createStatement();
        ResultSet rs = stmt.executeQuery("select last_build_id from bayzr_last_check where checker='sonarqube'");
        int last_build_id = 0;

        while(rs.next()){
           last_build_id  = rs.getInt("last_build_id");
           LOGGER.info("Get last build id: " + last_build_id);
           break;
        }
        rs.close();
        if (last_build_id==0){
           stmt.executeUpdate("insert into bayzr_last_check(checker, last_build_id) (select 'sonarqube', max(id) from bayzr_build_info where completed = 1)");
        } else {
           stmt.executeUpdate("update bayzr_last_check set last_build_id=(select max(id) from bayzr_build_info where completed = 1) where checker='sonarqube'");
        }
        String str_build = Integer.toString(last_build_id);
        rs = stmt.executeQuery("select bayzr_err, severity, file, pos, descript from bayzr_err_list where build_number = (select max(id) from bayzr_build_info where completed = 1 and id >= " + str_build + ") order by file, pos");
        while(rs.next()){
           String sev = rs.getString("severity");
           String file = rs.getString("file");
           String desc = rs.getString("descript");
           int err_tp = rs.getInt("bayzr_err");
           int pos = rs.getInt("pos");
           String err_rule = "";
           if sev=="" {
             if err_tp==2 {
                err_rule = "BayZRRule_High";
             } else if err_tp==1 {
                err_rule = "BayZRRule_Medium";
             } else {
                err_rule = "BayZRRule_Low"
             }
           } else {
                //err_rule = sev;
                err_rule = "BayZRRule_Low";
           }
           BayzrError dbError = new BayzrError(err_rule, desc, file, pos);
           fndIssuesList.add(dbError)
        }
        rs.close();
        stmt.close();
        conn.close();
      }catch(SQLException se){
        se.printStackTrace();
      }catch(Exception e){
        e.printStackTrace();
      }finally{
        try{
         if(stmt!=null)
            stmt.close();
        }catch(SQLException se2){
        }
        try{
         if(conn!=null)
            conn.close();
        }catch(SQLException se){
         se.printStackTrace();
        }
      }

      return fndIssuesList;
    }
  }

}
