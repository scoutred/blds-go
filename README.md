# blds-go

An API client for interacting with Building and Land Development permit APIs

The requirements of the spec are [here](https://github.com/open-data-standards/permitdata.org/wiki/Core-Permits-Dataset-Requirements)
The [parent page](https://github.com/open-data-standards/permitdata.org/wiki)
has optional portions of the dataset.

## issues

* Some of the dates are in `M/D/YYYY` despite the spec calling for `YYYY-MM-DD`
* Theres a `O` (letter, Oscar) in a numeric field
    * this comman can be used to solve it: `sed -e 's/,O,/,0,' -i -f data.csv`
* integers are encoded as strings in the json responses
