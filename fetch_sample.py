#!/usr/bin/python
import os
import urllib
import urllib2
import xml.etree.ElementTree as ET
from xml.dom.minidom import parseString

# API constants
API_KEY = '312e9d63e29aac03869e50be1ed6997b'
API_SECRET = 'b3516429fbcf0a69'
API_ENTRY = 'http://api.flickr.com/services/rest/'
API = 'flickr.photos.search'

# URL endpoint
BASE_URL = API_ENTRY + '?method=' + API + '&api_key=' + API_KEY

# search parameters
search_term_tags = ['black', 'white', 'red', 'orange', 'yellow', 'green', 'blue', 'purple', 'brown', 'grey']
search_sort = 'relevance'
search_per_page = 10

# my local path, yours would be different
local_path = '/Users/kevinnagaoka/Development/Mosaic/photos/'

# saves the images to disk
def save_images(color, img_urls):
	i = 1
	for url in img_urls:
		img_path = local_path + color
		# create directory to store colored images in, if it doesn't exist
		if not os.path.exists(img_path): 
			os.makedirs(img_path)
		# retrieve images from image URL and save to color directory
		urllib.urlretrieve(url, img_path + '/' + color + str(i) + '.jpg')
		i += 1

# gets the URL for each image
def get_image_urls(photo_list):
	photo_urls = []
	for photo in photo_list:
		attribs = photo.attrib
		farm_id = attribs['farm']
		server_id = attribs['server']
		id = attribs['id']
		secret = attribs['secret']
		img_src = 'http://farm' + farm_id + '.staticflickr.com/' + server_id + '/' + id + '_' + secret + '_s.jpg'
		photo_urls.append(img_src)
	return photo_urls

# query flickr with each color tag
for color in search_term_tags:
	query = BASE_URL + '&tags=' + color + '&sort=' + search_sort + '&per_page=' + str(search_per_page)
	file = urllib2.urlopen(query)
	tree = ET.parse(file)
	file.close()
	photos = tree.findall('photos/photo')
	photo_urls = get_image_urls(photos)
	save_images(color, photo_urls)
