#TODO
* Fix test onRubbishRequest


![Scheme](logo.png)

IGC Viewer is an online service that allows users to browse information about IGC files.

## Author
Jan Zimmer (janzim, 493594)

### Link to heroku
[http://igcviewer.herokuapp.com/igcinfo](http://igcviewer.herokuapp.com/igcinfo)

## Features
* Get meta information about the API

`GET /api`

* Track registration

`POST /api/igc`

* Returns the array of all tracks ids

`GET /api/igc`

* Returns the meta information about a given track with the provided id, or NOT FOUND response code with an empty body.

`GET /api/igc/<id>`

* Returns the single detailed meta information about a given track with the provided id, or NOT FOUND response code with an empty body. The response should always be a string, with the exception of the calculated track length, that should be a number.

`GET /api/igc/<id>/<field>`
