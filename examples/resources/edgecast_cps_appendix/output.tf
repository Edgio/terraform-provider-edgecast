 output "fetched_info_countrycode_Bermuda" {
    
    value = {
      for output in data.edgecast_cps_countrycodes.countrycodes.items : output.country =>
      output.two_letter_code if contains(["Bermuda"], output.country)
    }
  }

  output "fetched_info_countrycode_all" {
    value = data.edgecast_cps_countrycodes.countrycodes.items
  }

  output "fetched_info_countrycode_none" {
    value = {
      for output in data.edgecast_cps_countrycodes.countrycodes.items : output.country =>
      output.two_letter_code if contains(["random"], output.country)
    }
  }