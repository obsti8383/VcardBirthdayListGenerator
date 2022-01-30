# Introduction

![Go Vet and Lint Status](https://github.com/obsti8383/VcardBirthdayListGenerator/actions/workflows/go_lint_vet_and_testBuild.yml/badge.svg)

VcardBirthdayListGenerator generates a birthday list as csv or text (to stdout) from vcf files

Best to be used with vdirsyncer to download the VCARD files from a CardDav server
and the use VcardBirthdayListGenerator

References:
- vdirsyncer: https://github.com/pimutils/vdirsyncer 

# Example command lines
 ./VcardBirthdayListGenerator csv ~/.contacts/

 ./VcardBirthdayListGenerator text ~/.contacts/

