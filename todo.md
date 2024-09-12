# Todo
- RENOMMER TrkSummary, GpxSummary, car ne conviennent pas
    - []TrkSummary est en fait un résumé du Trk, qui est sur plusieurs points...

- modif calcTopo pour seulement faire le calcul
    - on mettre les valeurs cumules au trkpt qui possède un nom
    - on créera une func qui fera le calcul pour otenir info -detail


- modifier en CLI
    * gpx test
    * gpx trk 
        + info
        + ls
        + dist
    * gpx utils
        + dist : calculate distance from gpx to km 

- afficher carte + GPX 
    * trouver une API sur laquelle transmettre le GPX, et nous retourne la carte sur un URL ?  
      ----> https://www.visugpx.com/api/documentation_api.php 
    * Créer soi-même avec html + js  
      ----> https://wiki.openstreetmap.org/wiki/Openlayers_Track_example

- On peut ajouter des noms dans trkpt  
(trkpt est de type wptType, et contient donc `name`, `cmt`, `desc`)
On les utiliserait pour afficher les étapes de la rando


# Documentation gpx

https://www.topografix.com/GPX/1/1/

