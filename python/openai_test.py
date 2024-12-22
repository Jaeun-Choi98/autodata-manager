from dotenv import load_dotenv
import os
from openai import OpenAI

def serve():
  # .env 파일에서 환경 변수 로드
  load_dotenv()

  client = OpenAI(
    api_key = os.getenv("OPENAI_API_KEY")
  )

  table_data = '''
	OrderID, CustomerName, Product, Price, City
	1, Alice, Book, 10, NY
	2, Bob, Pen, 5, LA
	3, Alice, Book, 10, NY
	4, Charlie, Notebook, 15, SF
	'''
  # Create OpenAI API request
  request_data = {
    "model": "gpt-4o-mini",
    "messages": [
      {"role": "system", "content": "You are a database expert who can normalize database tables into 1NF, 2NF, and 3NF."},
      #{"role": "user", "content": f"Normalize the following table data to 1NF, 2NF, and 3NF. Provide the resulting table schemas in the format 'tablename: {{column1, column2, ...}}' for each normalization step:\n{table_data}"}
      ],
    "max_tokens": 500,
    "temperature": 0.7,
  }

  # Call OpenAI API
  response = client.chat.completions.create(**request_data)
  normalized_text = response.choices[0].message.content
  print(normalized_text)

if __name__ == '__main__':
  serve()