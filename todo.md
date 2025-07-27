# Todo

- features
    * calculation info between two nammed point
        * get intersection between trk
        * ...
    * get closest actual track using OSRM match function
    - **add func to get Ele, if missing**
        * https://geoservices.ign.fr/documentation/services/services-geoplateforme/altimetrie#72671
        * https://data.geopf.fr/altimetrie/swagger-ui/index.html#/Resources/get_resources_1_0_resources__get


- create function
    * check distance between trk[t].trkpt[-1] and trk[t+1].trkp[0]
        + try reverse (4 cases). If better, keep it
    * info -a : shows names in trkpt
    * info with rolling mean option
    * merge ls / info

- fix 
    * info -d: when last trkpt is named, there is still a `xxx --> end` entry

- refacto
    * make gpx.Filepath unexported (need SetFilepath, GetFilepath instead), to prevent save
    - gpx.Save : ajouter un calcul avant save pour obtenir les infos pour set denivPos, ...
    - modifier en CLI
        * gpx
            * apply
                + color
                + name
                + elevation
                + reverse
            * trk 
                + info
                + ls --> TODO: remove !
            * calc
                + dist : calculate distance from gpx to km 
                + calc-effort
            * utils
                + split
                + merge
                    + trk
                    + gpx
            * plot
            * test

- add check when no args in
    * dist
    * calc_effort



## someday, maybe
- afficher carte + GPX 
    * trouver une API sur laquelle transmettre le GPX, et nous retourne la carte sur un URL ?  
      ----> https://www.visugpx.com/api/documentation_api.php 
    * Créer soi-même avec html + js  
      ----> https://wiki.openstreetmap.org/wiki/Openlayers_Track_example


# Documentation gpx

https://www.topografix.com/GPX/1/1/

