from PIL import Image

def getAvgRGB():
   image=Image.open("./black.png")
   image=image.convert('RGB')
   img=list(image.getdata())
   r = 0
   g = 0
   b = 0
   for x in img:
      r += x[0]
      g += x[1]
      b += x[2]
   length = len(img)
   return (r / length, g / length, b / length)
