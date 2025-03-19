package utils


func IsPitchStrike(plate_x, plate_z float64) bool {

  plate_x_left := -0.71
  plate_x_right := 0.71
  plate_z_top := 3.5
  plate_z_bottom := 1.5

  return plate_x >= plate_x_left && plate_x <= plate_x_right && plate_z >= plate_z_bottom && plate_z <= plate_z_top 


}


