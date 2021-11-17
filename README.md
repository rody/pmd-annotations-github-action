# About

Github Action to create annotations from a [PMD](https://pmd.github.io) report.

---

## Usage

``` yaml
name: Analyse Source Code with PMD
on: [push]

jobs:
  analysis:
    name: Analysis
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v2

      # Generate the PMD Report
      #
      # You can use any action or build step
      # which can generate a PMD report in JSON
      # format (i.e. Maven, Gradle, ....)
      - name: PMD
        uses: rody/pmd-github-action@main
        with:
          rulesets: 'rulesets/apex/quickstart.xml'
          reportfile: 'pmd-report.json'
          format: 'json'
          failOnViolation: 'false'
          
      # Create annotations from PMD report
      - name: Create PMD annotations
        uses: rody/pmd-annotations-github-action@main
        with:
          reportfile: 'pmd-report.json'
          # Treat all rules with priority 1-4 as errors
          min-error-priority: 4
          # If any violations is an error, mark the step as failed
          fail-on-error: true
```

## Input

See [action.yml](action.yml)

# License

The scripts and documentation in this project are released under the [MIT License](LICENSE)

# Contributions

Contributions are welcome!

## Code of Conduct

:wave: Be nice and respectful.
