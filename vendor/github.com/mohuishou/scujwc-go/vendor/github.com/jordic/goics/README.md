
### A go toolkit for decoding, encoding icalendar ics ical files

Alpha status

[![Build Status](https://travis-ci.org/jordic/goics.svg?branch=master)](https://travis-ci.org/jordic/goics)

After trying to decode some .ics files and review available go packages, I decided to start writing this pacakge.

First attempt was from a fixed structure, similar to that needed. Later, I started to investigate the format and discovered that it is a pain, and has many variants, depending on who implements it. For this reason I evolved it to a tookit for decode and encode the format.

**Check examples dir for user cases:**

[Demo app encoding an ical, using a mysql DB source](examples/mysqlsource)


Features implemented:

- Parsing and writing of vevent, not completly.. 
- No recursive events, 
- And no alarms inside events


Follows:
http://tools.ietf.org/html/rfc5545


TODO
--

Integrate testing from:
https://github.com/libical/libical


CHANGELOG:
--

- v00. First api traial
- v01. Api evolves to a ical toolkit
- v02. Added example of integration with a mysql db.

Thanks to:
Joe Shaw Reviewing first revision.

