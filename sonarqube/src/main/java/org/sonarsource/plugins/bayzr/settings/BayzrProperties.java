package org.sonarsource.plugins.bayzr.settings;

import static java.util.Arrays.asList;

import java.util.List;

import org.sonar.api.config.PropertyDefinition;

public class BayzrProperties {

  public static final String DBPARAM_KEY = "sonar.bayzr.url";
  public static final String DBUS_KEY = "sonar.bayzr.user";
  public static final String DBPS_KEY = "sonar.bayzr.pass";
  public static final String SUFFIXES_KEY = "sonar.bayzr.files";
  public static final String LASTID_KEY = "sonar.bayzr.lastid";
  public static final String CATEGORY = "bayzr";

  private BayzrProperties() {
    // only statics
  }

  public static List<PropertyDefinition> definitions() {
    return asList(
      PropertyDefinition.builder(SUFFIXES_KEY)
        .name("List of checked files")
        .description("List of checked files")
        .defaultValue("cxx,.cpp,.cc,.c,.hxx,.hpp,.hh,.h")
        .category(CATEGORY)
        .build(),
      PropertyDefinition.builder(DBUS_KEY)
        .name("User")
        .description("Db connevtion parameters")
        .defaultValue("bayzr")
        .category(CATEGORY)
        .build(),
      PropertyDefinition.builder(DBPS_KEY)
        .name("Password")
        .description("Db connevtion parameters")
        .defaultValue("bayzr")
        .category(CATEGORY)
        .build(),
      PropertyDefinition.builder(LASTID_KEY)
        .name("Last build ID check")
        .description("Add build by ID more then last id")
        .defaultValue("no")
        .category(CATEGORY)
        .build()
      PropertyDefinition.builder(DBPARAM_KEY)
        .name("JDBC URL")
        .description("Db connevtion parameters (URL)")
        .defaultValue("jdbc:mysql://localhost/bayzr")
        .category(CATEGORY)
        .build()
      );
  }
}
