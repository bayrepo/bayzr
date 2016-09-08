package org.sonarsource.plugins.bayzr.rules;

import java.io.InputStream;
import java.nio.charset.StandardCharsets;
import java.util.Locale;

import org.sonar.api.server.rule.RulesDefinition;
import org.sonar.api.server.rule.RulesDefinitionXmlLoader;
import org.sonarsource.plugins.bayzr.languages.CppLanguage;

public final class BayzrRulesDefinition implements RulesDefinition {

  private static final String PATH_TO_RULES_XML = "/bayzr/cppbayzr-rules.xml";
  private static final String PATH_TO_RULES_XML_2 = "/bayzr/bayzr-compiler-gcc-rules.xml";
  private static final String PATH_TO_RULES_XML_3 = "/bayzr/bayzr-cppcheck-rules.xml";
  private static final String PATH_TO_RULES_XML_4 = "/bayzr/bayzr-rats-rules.xml";

  protected static final String KEY = "bayzr";
  protected static final String NAME = "BayZR";

  protected String rulesDefinitionFilePath() {
    return PATH_TO_RULES_XML;
  }

  protected String rulesDefinitionFilePat_2h() {
    return PATH_TO_RULES_XML_2;
  }

  protected String rulesDefinitionFilePath_3() {
    return PATH_TO_RULES_XML_3;
  }

  protected String rulesDefinitionFilePath_4() {
    return PATH_TO_RULES_XML_4;
  }

  private void defineRulesForLanguage(Context context, String repositoryKey, String repositoryName, String languageKey) {
    NewRepository repository = context.createRepository(repositoryKey, languageKey).setName(repositoryName);

    InputStream rulesXml = this.getClass().getResourceAsStream(rulesDefinitionFilePath());
    if (rulesXml != null) {
      RulesDefinitionXmlLoader rulesLoader = new RulesDefinitionXmlLoader();
      rulesLoader.load(repository, rulesXml, StandardCharsets.UTF_8.name());
    }

    InputStream rulesXml_2 = this.getClass().getResourceAsStream(rulesDefinitionFilePath_2());
    if (rulesXml_2 != null) {
      RulesDefinitionXmlLoader rulesLoader = new RulesDefinitionXmlLoader();
      rulesLoader.load(repository, rulesXml_2, StandardCharsets.UTF_8.name());
    }

    InputStream rulesXml_3 = this.getClass().getResourceAsStream(rulesDefinitionFilePath_3());
    if (rulesXml_3 != null) {
      RulesDefinitionXmlLoader rulesLoader = new RulesDefinitionXmlLoader();
      rulesLoader.load(repository, rulesXml_3, StandardCharsets.UTF_8.name());
    }

    InputStream rulesXml_4 = this.getClass().getResourceAsStream(rulesDefinitionFilePath_4());
    if (rulesXml_4 != null) {
      RulesDefinitionXmlLoader rulesLoader = new RulesDefinitionXmlLoader();
      rulesLoader.load(repository, rulesXml_4, StandardCharsets.UTF_8.name());
    }

    repository.done();
  }

  @Override
  public void define(Context context) {
    String repositoryKey = BayzrRulesDefinition.getRepositoryKeyForLanguage(CppLanguage.KEY);
    String repositoryName = BayzrRulesDefinition.getRepositoryNameForLanguage(CppLanguage.KEY);
    defineRulesForLanguage(context, repositoryKey, repositoryName, CppLanguage.KEY);
  }

  public static String getRepositoryKeyForLanguage(String languageKey) {
    return languageKey.toLowerCase(Locale.ENGLISH) + "-" + KEY;
  }

  public static String getRepositoryNameForLanguage(String languageKey) {
    return languageKey.toUpperCase(Locale.ENGLISH) + " " + NAME;
  }

}
