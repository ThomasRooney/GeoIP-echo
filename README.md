## GeoIP Echo server

I was looking for a way to tell which country my computer was operating in, and couldn't find anything that I could use from my terminal on the net (i.e., on a headless terminal).

IP echoing is a nice solution, and I have this alias set up in my .profile: alias ipecho='wget http://ipecho.net/plain -O - -q ; echo'

However what I was looking for required a bit more information, so I grabbed the GeoIP data from Maxmind, and used some go api bindings to read it for my purposes. This is the result:

Done pretty simply, though it was a bit of trial and error to get working on Heroku, as I didn't want to add the GeoIP database to this repository. Have a look in my web.go to see how it was done.

##TL;DR

alias geoipecho='wget http://ancient-shore-2349.herokuapp.com/ -O - -q; echo'

	$ geoipecho
	ip: 129.31.224.57
	Country: United Kingdom (GB)
	City: London
	Region: H9
	Postal Code: 
	Latitude: 51.514206
	Longitude: -0.093094


Feel free to use my code for whatever purpose you desire.
