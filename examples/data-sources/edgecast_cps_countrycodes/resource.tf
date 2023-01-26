data "edgecast_cps_countrycodes" "united_states" {
    name = "United States"
}

resource "edgecast_cps_certificate" "certificate_1" {
  country = data.edgecast_cps_countrycodes.united_states.items[0].two_letter_code
}