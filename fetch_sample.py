#!/usr/bin/python
import os
import urllib
import flickr

flickr.API_KEY = '312e9d63e29aac03869e50be1ed6997b'
flickr.API_SECRET = 'b3516429fbcf0a69'

# search parameters
search_term_tags = ['black', 'white', 'red', 'orange', 'yellow', 'green', 'blue', 'purple', 'brown', 'grey']
search_sort = 'relevance'
search_per_page = 10

# my local path, yours would be different
local_path = '/Users/kevinnagaoka/Development/Mosaic/photos/'

# loop through color tags
for color in search_term_tags:
	results = flickr.photos_search(tags = color, text = color, sort = search_sort, per_page = search_per_page)
	i = 1
	# go through list of results
	for photo in results:
		img_url = photo.getMedium()
		img_path = local_path + color
		# create directory to store colored images in, if it doesn't exist
		if not os.path.exists(img_path): 
			os.makedirs(img_path)
		# retrieve images from image URL and save to color directory
		urllib.urlretrieve(img_url, img_path + '/' + color + str(i) + '.jpg')
		i += 1
