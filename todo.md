# Todo
- [ ] modifier struct_gpx.go pour suivre la spec (eg. TrkSeg est un array)

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

