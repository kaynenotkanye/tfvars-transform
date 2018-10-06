package main

import (
  "os"
  "fmt"
  "flag"
  "bufio"
  "regexp"
  "strings"
)

func main() {
  // keyName must NOT be blank (bad things happen...)
  // I made it a requirement for keyName to be greater than 2 characters long.
  // 2 characters is an arbitrary number; typically, terraform variables should
  // hopefully be longer than 2 characters.
  tfvarsFile    := flag.String("f", "", "Filename of the tfvars file, example: terraform.tfvars")
  keyName       := flag.String("k", "", "Key name to assign/overwrite, example: appname_ami_id")
  keyValue      := flag.String("v", "", "Value to be populated, example: ami-abcd1234")

  // The following regexSuffix is aimed towards AWS AMI ID naming conventions
  // the primary purpose of this app is for a Packer AMI pipeline to update a tfvars file
  // in a downstream terraform pipeline; if the scope broadens, then the regex expression
  // will need to be changed to something more appropriate.
  regexSuffix    := flag.String("s", "[\\s]?=[\\s]?\"ami-[a-zA-Z0-9]{8,}\"", "Only use this flag when you need a custom non-AMI-related regex. Otherwise, do NOT use this flag.")
  flag.Parse()

  if len(*keyName) < 3 {
      fmt.Println("Program exiting. Please enter a key name greater than 2 characters.")
      os.Exit(1)
  }

  tempFile      := strings.Join([]string{*tfvarsFile, "tmp"}, ".")
  regexString   := strings.Join([]string{*keyName, *regexSuffix}, "")
  updatedString := strings.Join([]string{*keyName, " = \"", *keyValue, "\""}, "")  // need to escape the double quotes so that the ami_id is a tfvars string.

  var re = regexp.MustCompile(regexString)

  // Open the tfvars file as specified from the -f flag
  f, err := os.Open(*tfvarsFile)
  if err != nil {
    panic(err)
  }

  writeToFile, err := os.OpenFile(tempFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
  if err != nil {
    panic(err)
  }

  // Read from tfvars, replace, and write to temp file
  w := bufio.NewWriter(writeToFile)
  scanner := bufio.NewScanner(f) // f is the *os.File
  for scanner.Scan() {
      s := re.ReplaceAllString(scanner.Text(), updatedString)
      fmt.Fprintf(w, "%s\n", s)
  }
  if err := scanner.Err(); err != nil {
      fmt.Println("There was an error scanning the file: ", err)
      os.Exit(1)
  }
  f.Close()
  w.Flush()

  // Overwrite the tfvars file
  os.Rename(tempFile, *tfvarsFile)
  if err != nil {
    panic(err)
  }
}
