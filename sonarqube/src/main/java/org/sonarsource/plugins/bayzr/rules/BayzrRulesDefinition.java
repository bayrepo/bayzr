package org.sonarsource.plugins.bayzr.rules;

import java.io.InputStream;
import java.nio.charset.StandardCharsets;
import java.util.Locale;

import org.sonar.api.server.rule.RulesDefinition;
import org.sonar.api.server.rule.RulesDefinitionXmlLoader;
import org.sonarsource.plugins.bayzr.languages.CppLanguage;

public final class BayzrRulesDefinition implements RulesDefinition {

  private static final String PATH_TO_RULES_XML = "/bayzr/cppbayzr-rules.xml";

  protected static final String KEY = "bayzr";
  protected static final String NAME = "BayZR";

  protected String rulesDefinitionFilePath() {
    return PATH_TO_RULES_XML;
  }

  private void defineRulesForLanguage(Context context, String repositoryKey, String repositoryName, String languageKey) {
    NewRepository repository = context.createRepository(repositoryKey, languageKey).setName(repositoryName);

    InputStream rulesXml = this.getClass().getResourceAsStream(rulesDefinitionFilePath());
    if (rulesXml != null) {
      RulesDefinitionXmlLoader rulesLoader = new RulesDefinitionXmlLoader();
      rulesLoader.load(repository, rulesXml, StandardCharsets.UTF_8.name());
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
