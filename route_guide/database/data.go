package database

func GetNameData() []byte {
	var exampleData = []byte(`[{
    "location": {
        "latitude": 407838351,
        "longitude": -746143763
    },
    "name": "Patriots Path, Mendham, NJ 07945, USA"
},  {
    "location": {
        "latitude": 413843930,
        "longitude": -740501726
    },
    "name": "162 Merrill Road, Highland Mills, NY 10930, USA"
}, {
    "location": {
        "latitude": 406337092,
        "longitude": -740122226
    },
    "name": "6324 8th Avenue, Brooklyn, NY 11220, USA"
}, {
    "location": {
        "latitude": 406421967,
        "longitude": -747727624
    },
    "name": "1 Merck Access Road, Whitehouse Station, NJ 08889, USA"
}, {
    "location": {
        "latitude": 410248224,
        "longitude": -747127767
    },
    "name": "3 Hasta Way, Newton, NJ 07860, USA"
}]`)

	return exampleData
}

func GetAreaData() []byte {
	var exampleData = []byte(`[{
    "location": {
        "latitude": 407838351,
        "longitude": -746143763
    },
    "area": "Hilly"
},  {
    "location": {
        "latitude": 413843930,
        "longitude": -740501726
    },
    "area": "Snowy"
}, {
    "location": {
        "latitude": 406337092,
        "longitude": -740122226
    },
    "area": "Plateau"
}, {
    "location": {
        "latitude": 406421967,
        "longitude": -747727624
    },
    "area": "Plains"
}, {
    "location": {
        "latitude": 410248224,
        "longitude": -747127767
    },
    "area": "Coastal"
}]`)
	return exampleData
}

func GetLocData() []byte {
	var exampleData = []byte(`[
{
    "latitude": 407838351,
    "longitude": -746143763
  },
  {
    "latitude": 413843930,
    "longitude": -740501726
  },
  {
    "latitude": 406337092,
    "longitude": -740122226
  },
  {
    "latitude": 406421967,
    "longitude": -747727624
  },
  {
    "latitude": 410248224,
    "longitude": -747127767
  }
]`)
	return exampleData
}
