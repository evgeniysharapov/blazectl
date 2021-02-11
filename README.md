[![Build Status](https://travis-ci.com/samply/blazectl.svg?branch=master)](https://travis-ci.com/samply/blazectl)
[![Go Report Card](https://goreportcard.com/badge/github.com/samply/blazectl)](https://goreportcard.com/report/github.com/samply/blazectl)

# blazectl

blazectl is a command line tool to control your FHIR® server. blazectl also works with [Blaze][4].

Currently you can upload transaction bundles from a directory and count resources.

## Installation

blazectl is written in Go. All you need is a single binary which is available for Linux, macOS and Windows.

### Linux

1. Download the latest release with the command:

   ```bash
   curl -LO https://github.com/samply/blazectl/releases/download/v0.4.0/blazectl-0.4.0-linux-amd64.tar.gz
   ```

1. Untar the binary:

   ```bash
   tar xzf blazectl-0.4.0-linux-amd64.tar.gz
   ```
   
1. Move the binary in to your PATH.

   ```bash
   sudo mv ./blazectl /usr/local/bin/blazectl
   ```

1. Test to ensure the version you installed is up-to-date:

   ```bash
   blazectl --version
   ```

### macOS

1. Download the latest release with the command:

   ```bash
   curl -LO https://github.com/samply/blazectl/releases/download/v0.4.0/blazectl-0.4.0-darwin-amd64.tar.gz
   ```

1. Untar the binary:

   ```bash
   tar xzf blazectl-0.4.0-darwin-amd64.tar.gz
   ```
   
1. Move the binary in to your PATH.

   ```bash
   sudo mv ./blazectl /usr/local/bin/blazectl
   ```

1. Test to ensure the version you installed is up-to-date:

   ```bash
   blazectl --version
   ```

### Windows

1. Download the latest release [here][3]

1. Unzip the binary.

1. Add the binary in to your PATH.

1. Test to ensure the version you downloaded is up-to-date:

   ```
   blazectl --version
   ```
   
## Usage

```
$ blazectl
Usage:
  blazectl [command]

Available Commands:
  count-resources Counts all resources by type
  download        Download FHIR resources into an NDJSON file
  help            Help about any command
  upload          Upload transaction bundles

Flags:
  -h, --help            help for blazectl
      --password string password information for basic authentication
      --server string   the base URL of the server to use
      --user string     user information for basic authentication
      --version         version for blazectl

Use "blazectl [command] --help" for more information about a command.
```

### Upload

You can use the upload command to upload transaction bundles to our server. currently only JSON files are supported. If you don't have any transaction bundles, you can generate some with [SyntheaTM][5].

Assuming the URL of your FHIR server is `http://localhost:8080`, in order to upload run:

```bash
blazectl --server http://localhost:8080 upload my/bundles
```

You will see a progress bar with an estimated ETA during upload. After the upload, a statistic inspired by [vegeta][6] will be printed:

```
Starting Upload to http://localhost:8080 ...
Uploads          [total, concurrency]     362, 4
Success          [ratio]                  100 %
Duration         [total]                  1m42s
Requ. Latencies  [mean, 50, 95, 99, max]  826ms, 534ms, 2.71s, 3.85s 6.467s
Proc. Latencies  [mean, 50, 95, 99, max]  710ms, 526ms, 2.041s, 2.739s 4.133s
Bytes In         [total, mean]            5.10 MiB, 14.59 KiB
Bytes Out        [total, mean]            61.74 MiB, 176.59 KiB
Status Codes     [code:count]             200:362
```

The statistics have the following meaning:

* Uploads - the total number of files uploaded with the given concurrency
* Success - the success rate (possible errors will be printed under the statistics)
* Duration - the total duration of the upload
* Requ. Latencies - mean, max and percentiles of the duration of whole requests including networks transfers 
* Proc. Latencies - mean, max and percentiles of the duration of the server processing time excluding networks transfers 
* Bytes In - total and mean number of bytes returned by the server
* Bytes Out - total and mean number of bytes send by blazectl
* Status Codes - a list of status code frequencies. Will show non-200 status codes if they happen.

### Download

You can use the download command to download bundles from the server. Downloaded bundles are stored within an NDJSON file. This operation is non-destructive on your site, i.e. if the specified NDJSON file already exists then it won't be overwritten.

Use the download command as follows:

```bash
blazectl --server http://localhost:8080/fhir download
         --type Patient
         --query "gender:not=male&_count=10"
         --output-file "~/Downloads/Patients.ndjson"
```

Next to the mandatory FHIR resource type you can also optionally specify a valid FHIR search query to limit downloaded bundles. The query must not start with a `?` token.


As soon as the download has finished you will be shown a download statistics overview that looks something like this:

```
Pages           [total]                 184
Resources       [total]                 1835
Resources/Page  [min, mean, max]        5, 9, 10
Duration        [total]                 371ms
Requ. Latencies	[mean, 50, 95, 99, max]	1ms, 1ms, 2ms, 2ms, 3ms
Proc. Latencies	[mean, 50, 95, 99, max]	1ms, 1ms, 1ms, 2ms, 3ms
Bytes In        [total, mean]           1.22 MiB, 6.82 KiB
```

The statistics have the following meaning:

* Pages - total number of pages requested from the server to retrieve resources
* Resources - total number of downloaded resources
* Resources/Page - minimum, mean and maximum number of resources over all pages 
* Duration - total duration of the download
* Requ. Latencies - mean, max and percentiles of the duration of whole requests including networks transfers
* Proc. Latencies - mean, max and percentiles of the duration of the server processing time excluding network transfers
* Bytes In - total and mean number of bytes returned by the server

### Count Resources

The count-resources command is useful to see how many resources a FHIR server stores by resource type. The resource counting is done by first fetching the capability statement of the server. After that blazectl will perform a search-type interaction with query parameter `_summary` set to `count` on every resource type which supports that interaction using one batch request. Bundle.total will be used as resource count.

You can run:
 
```bash
blazectl --server http://localhost:8080 count-resources
```

It will return:

```
Count all resources on http://localhost:8080 ...

AllergyIntolerance       :    7297
CarePlan                 :   49818
Claim                    :  689111
Condition                :  116688
DiagnosticReport         :  193141
Encounter                :  540542
ExplanationOfBenefit     :  540542
Goal                     :   39857
ImagingStudy             :   11212
Immunization             :  187987
MedicationAdministration :    6400
MedicationRequest        :  148569
Observation              : 2689215
Organization             :   52645
Patient                  :   16875
Practitioner             :   52647
Procedure                :  418310
```

## Similar Software

* [VonkLoader][1] - can also upload transaction bundles but needs .NET SDK
* [Synthea Uploader][2] - no parallel uploads

## License

Copyright 2019 The Samply Development Community

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.

[1]: <http://docs.simplifier.net/vonkloader/>
[2]: <https://github.com/synthetichealth/uploader>
[3]: <https://github.com/samply/blazectl/releases/download/v0.4.0/blazectl-0.4.0-windows-amd64.zip>
[4]: <https://github.com/samply/blaze>
[5]: <https://github.com/synthetichealth/synthea>
[6]: <https://github.com/tsenart/vegeta>
