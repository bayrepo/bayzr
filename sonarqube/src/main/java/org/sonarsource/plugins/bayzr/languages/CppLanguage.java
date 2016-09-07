package org.sonarsource.plugins.bayzr.languages;

import java.util.ArrayList;
import java.util.List;

import org.apache.commons.lang.StringUtils;
import org.sonar.api.config.Settings;
import org.sonar.api.resources.AbstractLanguage;

/**
 * This class defines the fictive Foo language.
 */
public final class CppLanguage extends AbstractLanguage {

  public static final String NAME = "Cpp";
  public static final String KEY = "cpp";
  //public static final String FILE_SUFFIXES_PROPERTY_KEY = "sonar.cpp.file.suffixes";
  public static final String FILE_SUFFIXES_PROPERTY_KEY = "sonar.bayzr.files";
  public static final String DEFAULT_FILE_SUFFIXES = "cxx,.cpp,.cc,.c,.hxx,.hpp,.hh,.h";

  private final Settings settings;

  public CppLanguage(Settings settings) {
    super(KEY, NAME);
    this.settings = settings;
  }

  @Override
  public String[] getFileSuffixes() {
    String[] suffixes = filterEmptyStrings(settings.getStringArray(FILE_SUFFIXES_PROPERTY_KEY));
    if (suffixes.length == 0) {
      suffixes = StringUtils.split(DEFAULT_FILE_SUFFIXES, ",");
    }
    return suffixes;
  }

  private String[] filterEmptyStrings(String[] stringArray) {
    List<String> nonEmptyStrings = new ArrayList<>();
    for (String string : stringArray) {
      if (StringUtils.isNotBlank(string.trim())) {
        nonEmptyStrings.add(string.trim());
      }
    }
    return nonEmptyStrings.toArray(new String[nonEmptyStrings.size()]);
  }

}
