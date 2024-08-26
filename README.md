# Flow Log Tagging Processor

## Overview

This tool processes AWS flow logs to count tag occurrences and generate number of port/protocol combinations.

## Assumptions

- **Log Format**: The program only supports the default AWS flow version 2 log format
- **Protocol Mapping**: This program uses a file `protocol-numbers-1.csv` to map protocol numbers to text. The file was taken from AWS.
- **File Sizes**: The program is optimized for processing flow log files up to 10 MB and lookup tables with up to 10,000 entries.
- **Concurrency**: The program uses concurrency to process the flow log file in chunks, which improves performance for larger files.

## Installation

### Prerequisites

- **Go 1.16+** installed on your system.

### Instructions

1. **Clone the repository:**
    ```bash
    git clone https://github.com/prateek1192/flow-log-tagging.git
    cd flow-log-tagging
    ```

2. **Build the project:**
    ```bash
    make build
    ```

3. **Run the program:**
    ```bash
    ./flowlog-processor --lookupfile=mylookup.csv --flowlogfile=myflowlog.txt --outputfile=myoutput.csv
    ```

   If no arguments are provided, default files will be used:
   - `lookup_table.csv` for the lookup table
   - `flowLog.txt` for the flow log file
   - `tag_counts.csv` for the output file

### Makefile Commands

- **Build the project:**
    ```bash
    make build
    ```
## Functionality

### Command-Line Arguments

- `--lookupfile` : Path to the CSV file containing the lookup table. Default is `lookup_table.csv`.
- `--flowlogfile` : Path to the file containing flow logs. Default is `flowLog.txt`.
- `--outputfile` : Path to the output CSV file where tag counts and port/protocol statistics will be saved. Default is `tag_counts.csv`.

### Example Usage

To run the processor with custom file paths:
```bash
./flowlog-processor --lookupfile=lookup.csv --flowlogfile=flowlog.txt --outputfile=output.csv
