# Introduction
VcardBirthdayListGenerator generates a birthday list as csv (to stdout) from vcf files

Best to be used with vdirsyncer to download the VCARD files from a CardDav server
and the use VcardBirthdayListGenerator

References:
- vdirsyncer: https://github.com/pimutils/vdirsyncer 

# Command line parameters
The following parameters are available:

  --path : path where the vcf files reside (or vcf file directly) (required)

The following commands are available:
		
  version : prints version

# Example command lines
./VcardBirthdayListGenerator --path ~/.contacts/

./VcardBirthdayListGenerator version
