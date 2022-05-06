# bhg-final-project
This project is designed to test how files can be exfiltrated from
a target machine. The program will locate all files matching a given
regular expression pattern and create an HTTP POST with the data and
the filenames.

# Usage
Call the fileSearch executable. The filesystem directory to search
and the destination URL are hardcoded so that no user input is needed.
./fileSearch 
