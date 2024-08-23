# Usage

```bash
gpx-cli ls      # print all trk name
gpx-cli ls -a  # include trkpt names


# calculate distance effort for each trk
gpx-cli info
gpx-cli info false # disable ascii format

gpx-cli info-detail    # wip: distance of each step on a trk

# Reverse trk and output to out/toto.xml
gpx-cli reverse all # reverse all trk
gpx-cli reverse 5   # reverse 5-ieth trk

# utils

# distance between 2 points in meters
gpx-cli dist lat1 lon1 lat2 lon2

gpx-cli calc_effort km denivPos denivNeg


```