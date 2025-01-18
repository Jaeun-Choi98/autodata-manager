import cv2
import pytesseract
import matplotlib.pyplot as plt

pytesseract.pytesseract.tesseract_cmd = r'C:\Program Files\Tesseract-OCR\tesseract.exe'

def plt_imshow(title='image', img=None, figsize=(8 ,5)):
  plt.figure(figsize=figsize)

  if type(img) == list:
      if type(title) == list:
          titles = title
      else:
          titles = []

          for i in range(len(img)):
              titles.append(title)

      for i in range(len(img)):
          if len(img[i].shape) <= 2:
              rgbImg = cv2.cvtColor(img[i], cv2.COLOR_GRAY2RGB)
          else:
              rgbImg = cv2.cvtColor(img[i], cv2.COLOR_BGR2RGB)

          plt.subplot(1, len(img), i + 1), plt.imshow(rgbImg)
          plt.title(titles[i])
          plt.xticks([]), plt.yticks([])

      plt.show()
      
  else:
      if len(img.shape) < 3:
          rgbImg = cv2.cvtColor(img, cv2.COLOR_GRAY2RGB)
      else:
          rgbImg = cv2.cvtColor(img, cv2.COLOR_BGR2RGB)

      plt.imshow(rgbImg)
      plt.title(title)
      plt.xticks([]), plt.yticks([])
      plt.show()

# 이미지 로드
image = cv2.imread('../test.png')

gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY) 
thresh = cv2.threshold(gray, 127, 255, cv2.THRESH_BINARY + cv2.THRESH_OTSU)[1]

# 텍스트 추출
text = pytesseract.image_to_string(thresh)

# 결과 출력
print(text)

plt_imshow("",[thresh,gray,image])




