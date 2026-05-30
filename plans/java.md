# Java implementation plan

**Read `plans/shared-context.md` first.** This document covers only Java-specific decisions.

---

## Goal

A Java CLI application in `java/` that reads `resume.yaml`, groups work entries by employer, and writes `docs/index.html`. Uses modern Java idioms: records, sealed interfaces, pattern matching, and text blocks.

---

## Language version

- **Java 26** (latest GA release). Set `<java.version>26</java.version>`. Use `source`/`target` 26 in the Maven compiler plugin. Java 25 is the current LTS; either is acceptable, but the plans target 26 for latest-stable alignment.
- Records, pattern matching, and sealed interfaces are all stable (no `--enable-preview` needed).
- No module-info.java required; a simple unnamed module is fine.

---

## Build system

**Maven** (`pom.xml`). Produces an executable fat-jar via `maven-shade-plugin`.

```xml
<groupId>io.github.stephenbrown2</groupId>
<artifactId>resume-renderer</artifactId>
<version>1.0-SNAPSHOT</version>
<packaging>jar</packaging>

<properties>
  <java.version>26</java.version>
  <maven.compiler.source>26</maven.compiler.source>
  <maven.compiler.target>26</maven.compiler.target>
</properties>
```

---

## Dependencies (`pom.xml`)

```xml
<!-- YAML parsing -->
<dependency>
  <groupId>com.fasterxml.jackson.dataformat</groupId>
  <artifactId>jackson-dataformat-yaml</artifactId>
  <version>2.18.x</version>  <!-- use latest 2.18.x; Jackson 3.x requires Java 17+ and is used by networknt 3.x -->
</dependency>
<dependency>
  <groupId>com.fasterxml.jackson.core</groupId>
  <artifactId>jackson-databind</artifactId>
  <version>2.18.x</version>
</dependency>

<!-- JSON Schema validation (Draft 2020-12) -->
<dependency>
  <groupId>com.networknt</groupId>
  <artifactId>json-schema-validator</artifactId>
  <version>3.0.3</version>  <!-- requires Java 17+; use Jackson 3.x -->
</dependency>

<!-- Templating -->
<dependency>
  <groupId>io.pebbletemplates</groupId>
  <artifactId>pebble</artifactId>
  <version>3.2.x</version>  <!-- use latest 3.2.x -->
</dependency>
```

**Pebble** is a Jinja2/Twig-style templating engine with automatic HTML escaping, `for` loops with `loop.last`, `|` filters, and macros. It is well-maintained, has no transitive surprise dependencies, and is significantly simpler than Freemarker for this use case.

**networknt json-schema-validator 3.x** uses Jackson 3.x internally. If Jackson version conflicts arise between 2.x (YAML) and 3.x (schema-validator), resolve by using Jackson 3.x for both - `jackson-dataformat-yaml` has a 3.x line compatible with Java 17+. Alternatively, use version 2.x of the schema-validator with Jackson 2.x to avoid the conflict entirely.

---

## File structure

```
java/
  pom.xml
  src/
    main/
      java/
        io/github/stephenbrown2/resume/
          Main.java            # entry point, CLI args, orchestration
          model/
            Resume.java        # top-level record
            Basics.java
            WorkEntry.java
            Skills.java
            SkillSet.java
            SkillItem.java
            Project.java
            Certificate.java
            Education.java
            Language.java
            Interest.java
            Testimonial.java
            Reference.java
            Disposition.java
            Relocation.java
            Location.java
            Profile.java
          grouping/
            EmployerGroup.java  # record for grouped work
            WorkGrouper.java    # static grouping logic
          render/
            HtmlRenderer.java   # Pebble setup and rendering
            DateFormatter.java  # ISO date → display string
            NbspFilter.java     # custom Pebble filter
            LevelClassFilter.java
      resources/
        template.html          # Pebble template
```

---

## Data model

Use **Java records** for all model types. Jackson deserializes into records via `@JsonProperty` annotations or a constructor-based approach (Jackson 2.18 supports record deserialization natively with `jackson-module-parameter-names` or via `@JsonCreator`).

Configure the `ObjectMapper` with:
```java
ObjectMapper mapper = new ObjectMapper(new YAMLFactory());
mapper.configure(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES, false);
mapper.setPropertyNamingStrategy(PropertyNamingStrategies.LOWER_CAMEL_CASE);
```

Example records:

```java
public record Resume(
    Basics basics,
    Disposition disposition,
    List<WorkEntry> work,
    List<Project> projects,
    Skills skills,
    List<Certificate> certificates,
    List<Education> education,
    List<Language> languages,
    List<Interest> interests,
    List<Testimonial> testimonials,
    List<Reference> references
) {}

public record WorkEntry(
    String employer,
    @JsonProperty("employerGroup") String employerGroup,  // nullable
    String position,
    String url,
    @JsonProperty("startDate") String startDate,
    @JsonProperty("endDate") String endDate,              // nullable
    String summary,
    String location,
    List<String> highlights,
    List<String> keywords
) {}

public record Skills(
    List<SkillSet> sets,
    List<SkillItem> list
) {}

public record SkillSet(String name, List<String> skills) {}
public record SkillItem(String name, String level, String summary, Integer years) {}
```

Where a field may be absent in YAML, use `@JsonProperty(required = false)` or rely on `FAIL_ON_UNKNOWN_PROPERTIES = false` with null defaults.

---

## Employer grouping (`WorkGrouper.java`)

```java
public class WorkGrouper {
    public static List<EmployerGroup> group(List<WorkEntry> entries) {
        // Walk entries, compute key = employerGroup != null ? employerGroup : employer
        // Accumulate consecutive same-key entries into EmployerGroup
        // Track formerNames as ordered set of distinct employer strings
        // Compute group startDate (min) and endDate (max/null=Present)
    }
}

public record EmployerGroup(
    String displayName,
    List<String> formerNames,
    String url,
    String startDate,
    String endDate,           // null means "Present"
    List<WorkEntry> positions
) {}
```

---

## Date formatting (`DateFormatter.java`)

```java
public class DateFormatter {
    private static final String[] MONTHS =
        {"Jan","Feb","Mar","Apr","May","Jun","Jul","Aug","Sep","Oct","Nov","Dec"};

    public static String format(String iso) {
        if (iso == null || iso.isBlank()) return "Present";
        String[] parts = iso.split("-");
        // parts[0] = year, parts[1] = month (optional), parts[2] = day (optional)
        if (parts.length >= 2) {
            int month = Integer.parseInt(parts[1]) - 1;
            return MONTHS[month] + " " + parts[0];
        }
        return parts[0]; // year only
    }
}
```

---

## Non-breaking space filter (`NbspFilter.java`)

Implement `com.mitchellbosecke.pebble.extension.Filter`:

```java
public class NbspFilter implements Filter {
    @Override
    public Object apply(Object input, Map<String, Object> args, PebbleTemplate self,
                        EvaluationContext context, int lineNumber) {
        if (!(input instanceof String s)) return input;
        String[] words = s.split(" ", -1);
        StringBuilder sb = new StringBuilder();
        for (int i = 0; i < words.length; i++) {
            sb.append(words[i]);
            if (i < words.length - 1) {
                boolean shortWord = words[i].length() <= 4;
                boolean nextLonger = words[i + 1].length() > 4;
                sb.append(shortWord && nextLonger ? "&nbsp;" : " ");
            }
        }
        return new RawString(sb.toString()); // or use Pebble's raw output mechanism
    }
}
```

Register via a custom `Extension`:

```java
public class ResumeExtension extends AbstractExtension {
    @Override
    public Map<String, Filter> getFilters() {
        return Map.of(
            "nbsp_words",  new NbspFilter(),
            "format_date", new DateFormatFilter(),
            "level_class", new LevelClassFilter()
        );
    }
}
```

For raw (unescaped) filter output in Pebble, use `{% autoescape false %}...{% endautoescape %}` around the summary, or use Pebble's `raw` filter chained after `nbsp_words`.

---

## Level class filter (`LevelClassFilter.java`)

```java
// input: "Advanced" → "adv", "Intermediate" → "mid", anything else → ""
```

---

## Pebble template (`src/main/resources/template.html`)

Pebble syntax is nearly identical to Jinja2. Key differences from the Rust/Python plans:
- Loop variable is `loop.index` (1-based); use `loop.last` for "is last".
- Filters chain with `|`: `{{ date | format_date }}`.
- Macros: `{% macro job_tags(keywords) %}...{% endmacro %}` for DRY tag rendering.
- Use `{{ variable | raw }}` to output pre-escaped HTML from filters that produce `&nbsp;`.

Skills lookup - in Pebble, build a map in Java and pass it as `skillMap: Map<String, SkillItem>`:
```java
Map<String, SkillItem> skillMap = skills.list().stream()
    .collect(Collectors.toMap(SkillItem::name, s -> s));
```
Then in the template: `{{ skillMap[skillName].level | level_class }}`.

---

## Schema validation (`Main.java`)

After reading the YAML, convert it to a JSON string (via Jackson) and validate against `schema.json` using networknt's validator:

```java
import com.networknt.schema.*;
import com.fasterxml.jackson.databind.JsonNode;

static void validateSchema(JsonNode data, String schemaPath) throws Exception {
    JsonSchemaFactory factory = JsonSchemaFactory.getInstance(SpecVersion.VersionFlag.V202012);
    SchemaValidatorsConfig config = SchemaValidatorsConfig.builder().build();
    JsonSchema schema = factory.getSchema(
        URI.create("file://" + Path.of(schemaPath).toAbsolutePath()),
        config
    );
    Set<ValidationMessage> errors = schema.validate(data);
    if (!errors.isEmpty()) {
        errors.forEach(e -> System.err.println("validation error: " + e.getMessage()));
        System.exit(1);
    }
}
```

The YAML `ObjectMapper` produces a `JsonNode` tree that networknt can validate directly - no intermediate JSON string needed.

## `--name-font` flag and Google Fonts URL

Add to `Main.java`'s arg parsing:

```java
String nameFont = "Instrument Serif";
boolean skipValidation = false;
for (int i = 0; i < args.length; i++) {
    // ...existing flags...
    if (args[i].equals("--name-font") || args[i].equals("-f")) nameFont = args[++i];
    if (args[i].equals("--skip-validation"))                    skipValidation = true;
}
String fontUrl        = nameFont.replace(" ", "+");
String googleFontsLink = String.format(
    "<link href=\"https://fonts.googleapis.com/css2?family=%s:ital@0;1&amp;display=swap\" rel=\"stylesheet\">",
    fontUrl);
String nameFontCSS    = String.format("'%s', Georgia, serif", nameFont);
```

Pass `googleFontsLink` and `nameFontCSS` to `HtmlRenderer` and from there into the Pebble template context as raw strings. In Pebble, output them with `{{ googleFontsLink | raw }}` and `{{ nameFontCSS | raw }}` to prevent double-escaping.

## `Main.java`

```java
public class Main {
    public static void main(String[] args) throws Exception {
        // Parse --input / -i, --output / -o, --name-font / -f, --skip-validation
        String input    = "../resume.yaml";
        String output   = "../docs/index.html";
        String nameFont = "Instrument Serif";
        boolean skipVal = false;
        for (int i = 0; i < args.length; i++) {
            if (args[i].equals("--input")  || args[i].equals("-i")) input    = args[++i];
            if (args[i].equals("--output") || args[i].equals("-o")) output   = args[++i];
            if (args[i].equals("--name-font") || args[i].equals("-f")) nameFont = args[++i];
            if (args[i].equals("--skip-validation"))                   skipVal  = true;
        }

        ObjectMapper mapper = new ObjectMapper(new YAMLFactory())
            .configure(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES, false);
        JsonNode rawTree = mapper.readTree(Path.of(input).toFile());
        if (!skipVal) validateSchema(rawTree, "schema.json");
        Resume resume = mapper.treeToValue(rawTree, Resume.class);

        List<EmployerGroup> groups = WorkGrouper.group(resume.work());
        String html = HtmlRenderer.render(resume, groups);

        Files.writeString(Path.of(output), html, StandardOpenOption.CREATE,
                          StandardOpenOption.TRUNCATE_EXISTING);
        System.err.println("wrote " + output);
    }
}
```

---

## `java/README.md`

Create `java/README.md` documenting this implementation. It should cover:

- **Prerequisites:** Java 26+ JDK and Maven 3.9+. Install the JDK via [SDKMAN](https://sdkman.io/) (`sdk install java 26-open`) or a package manager; Maven via `sdk install maven` or `brew install maven`.
- **Build:** `mvn -q package -DskipTests` (produces `target/resume-renderer-1.0-SNAPSHOT.jar`; or `just java-build` from the repo root).
- **Run:** `java -jar target/resume-renderer-1.0-SNAPSHOT.jar [flags]` (or `just java-render` from the repo root).
- **Flags:** table matching the CLI interface in `shared-context.md` (`--input`, `--output`, `--name-font`, `--skip-validation`).
- **Output:** writes `docs/index.html` (relative to the repo root when using the default path).
- **Fat jar:** note that `maven-shade-plugin` bundles all dependencies - the jar is fully self-contained and needs only a JRE to run.

---

## Build and run

```sh
cd java
mvn -q package -DskipTests
java -jar target/resume-renderer-1.0-SNAPSHOT.jar \
     --input ../resume.yaml --output ../docs/index.html
```

Add to repo `justfile`:

```just
java-build:
    cd java && mvn -q package -DskipTests

java-render: java-build
    java -jar java/target/resume-renderer-1.0-SNAPSHOT.jar \
         --input resume.yaml --output docs/index.html
```

---

## Notes

- Jackson 2.18 deserializes Java records directly; no `@JsonCreator` boilerplate needed if field names match YAML keys (after camelCase mapping).
- Pebble 3.x auto-escapes HTML by default. Literal HTML entities in the template text (`&middot;`, `&amp;`) are passed through unchanged. Filter output that contains entities must use `| raw` or the `{% autoescape false %}` block.
- `List.copyOf()` and `Map.of()` (available since Java 9) make the grouping logic concise without Guava.
- Avoid streams where a simple `for` loop is clearer - this is a small transformation and readability matters.
- Do not use any preview features (`--enable-preview`) - records and pattern matching are stable in Java 21+.
