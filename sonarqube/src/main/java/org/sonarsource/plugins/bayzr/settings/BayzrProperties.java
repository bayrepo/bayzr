package org.sonarsource.plugins.bayzr.settings;

import static java.util.Arrays.asList;

import java.util.List;

import org.sonar.api.config.PropertyDefinition;

public class BayzrProperties {

  public static final String DB_KEY = "sonar.bayzr.jdbcdriver";
  public static final String DBPARAM_KEY = "sonar.bayzr.url";
  public static final String CATEGORY = "bayzr";

  private BayzrProperties() {
    // only statics
  }

  public static List<PropertyDefinition> definitions() {
    return asList(
      PropertyDefinition.builder(DB_KEY)
        .name("JDBC Driver")
        .description("Db connevtion parameters")
        .defaultValue(String.valueOf(false))
        .category(CATEGORY)
        .build(),
      PropertyDefinition.builder(DBPARAM_KEY)
        .name("JDBC URL")
        .description("Db connevtion parameters (URL)")
        .defaultValue(String.valueOf(false))
        .category(CATEGORY)
        .build()
      );
  }
}
