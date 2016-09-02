package org.sonarsource.plugins.bayzr;

import org.sonar.api.Plugin;
import org.sonarsource.plugins.bayzr.hooks.DisplayIssuesInScanner;
import org.sonarsource.plugins.bayzr.hooks.DisplayQualityGateStatus;
import org.sonarsource.plugins.bayzr.languages.CppLanguage;
import org.sonarsource.plugins.bayzr.languages.CppQualityProfile;
import org.sonarsource.plugins.bayzr.measures.ComputeSizeAverage;
import org.sonarsource.plugins.bayzr.measures.ComputeSizeRating;
import org.sonarsource.plugins.bayzr.measures.BayzrMetrics;
import org.sonarsource.plugins.bayzr.measures.SetSizeOnFilesSensor;
import org.sonarsource.plugins.bayzr.rules.BayzrIssuesLoaderSensor;
import org.sonarsource.plugins.bayzr.rules.BayzrRulesDefinition;
import org.sonarsource.plugins.bayzr.settings.BayzrProperties;
import org.sonarsource.plugins.bayzr.settings.SayHelloFromScanner;

/**
 * This class is the entry point for all extensions. It is referenced in pom.xml.
 */
public class BayzrPlugin implements Plugin {

  @Override
  public void define(Context context) {
    // tutorial on hooks
    // http://docs.sonarqube.org/display/DEV/Adding+Hooks
    context.addExtensions(DisplayIssuesInScanner.class, DisplayQualityGateStatus.class);

    // tutorial on languages
    context.addExtensions(CppLanguage.class, CppQualityProfile.class);

    // tutorial on measures
    context
      .addExtensions(BayzrMetrics.class, SetSizeOnFilesSensor.class, ComputeSizeAverage.class, ComputeSizeRating.class);

    // tutorial on rules
    context.addExtensions(BayzrRulesDefinition.class, BayzrIssuesLoaderSensor.class);

    // tutorial on settings
    context
      .addExtensions(BayzrProperties.definitions())
      .addExtension(SayHelloFromScanner.class);

  }
}
