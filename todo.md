# Todo

- create function
    * check distance between trk[t].trkpt[-1] and trk[t+1].trkp[0]
        + try reverse (4 cases). If better, keep it

    * info -a : shows names in trkpt
    * merge ls / info

- use rolling calcultion in info, etc. (improve elevation estimation)

- make gpx.Filepath unexported (need SetFilepath, GetFilepath instead), to prevent save

- **add func to get Ele, if missing**
    * https://geoservices.ign.fr/documentation/services/services-geoplateforme/altimetrie#72671
    * https://data.geopf.fr/altimetrie/swagger-ui/index.html#/Resources/get_resources_1_0_resources__get

- gpx.Save : ajouter un calcul avant save pour obtenir les infos pour set denivPos, ...

- add check when no args in
    * dist
    * calc_effort


- modifier en CLI
    * gpx test
    * gpx trk 
        + info
        + ls
        + dist
    * gpx utils
        + dist : calculate distance from gpx to km 

## someday, maybe
- afficher carte + GPX 
    * trouver une API sur laquelle transmettre le GPX, et nous retourne la carte sur un URL ?  
      ----> https://www.visugpx.com/api/documentation_api.php 
    * Créer soi-même avec html + js  
      ----> https://wiki.openstreetmap.org/wiki/Openlayers_Track_example


# Documentation gpx

https://www.topografix.com/GPX/1/1/

