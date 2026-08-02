# Todo

<<<<<<< HEAD
=======
- rename Trkpts to Section ?

- refacto :
    * simplify summary struct... how ?
        + or split gpx, trk, trpt package ? each containing summary

>>>>>>> refacto/reorg_export_stats
- features
    * get closest actual track using OSRM match function


- create function
    * check distance between trk[t].trkpt[-1] and trk[t+1].trkp[0]
        + try reverse (4 cases). If better, keep it
    * merge ls / info

- fix 
    * info -d: when last trkpt is named, there is still a `xxx --> end` entry

- add check when no args in
    * dist
    * calc_effort
    * 
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

