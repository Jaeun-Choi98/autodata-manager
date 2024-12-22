
from dotenv import load_dotenv
import os
import grpc
from concurrent import futures
from openai import OpenAI
import normalize_pb2_grpc
import normalize_pb2

load_dotenv()
client = OpenAI(
  api_key = os.getenv("OPENAI_API_KEY")
)

class NormalizeService(normalize_pb2_grpc.NormalizeServiceServicer):
  def NormalizeTable(self, request, context):
    table_data = request.data
    """
    request_data = {
    "model": "gpt-4o-mini",
    "messages": [
    {
        "role": "system",
        "content": "You are a database expert who can normalize database tables into 3NF." 
    },
    {
        "role": "user",
        "content": f"Normalize the following table data to 3NF. Provide the resulting table schemas in format csv '{{tableName1: [[\"col1\",\"col2\",\"col3\",..], [\"val1\",\"val2\",\"val3\",..],..], tableName2:[...]}}':'.\n\ndata:\n{table_data}"
    }],
    "max_tokens": 500,
    "temperature": 0.7,
    }
    response = client.chat.completions.create(**request_data)
    normalized_text = response.choices[0].message.content
    return normalize_pb2.NormalizedTable(result=normalized_text)
    """
    # oepnai를 사용했을 때 응답 값 e.g.
    """
    "choices": [
        {
          "message": {
            "role": "assistant",
            "content": "1NF:\nOrder: {OrderID, CustomerName, CustomerAddress, ProductName, ProductCategory, Quantity, Price, OrderDate}\n\n2NF:\nOrder: {OrderID, CustomerID, ProductID, Quantity, Price, OrderDate}\nCustomer: {CustomerID, CustomerName, CustomerAddress}\nProduct: {ProductID, ProductName, ProductCategory}\n\n3NF:\nOrder: {OrderID, CustomerID, ProductID, Quantity, Price, OrderDate}\nCustomer: {CustomerID, CustomerName, CustomerAddress}\nProduct: {ProductID, ProductName, ProductCategory}\n"
          }
        },...]
    """
    ret = "{\n  \"Customers\": [\n    [\"CustomerID\", \"CustomerName\", \"CustomerAddress\"],\n    [\"1\", \"Alice Kim\", \"Seoul, South Korea\"],\n    [\"2\", \"Bob Lee\", \"Busan, South Korea\"],\n    [\"3\", \"Charlie Park\", \"Daegu, South Korea\"],\n    [\"4\", \"David Choi\", \"Incheon, South Korea\"]\n  ],\n  \"Products\": [\n    [\"ProductID\", \"ProductName\", \"ProductCategory\"],\n    [\"1\", \"Laptop\", \"Electronics\"],\n    [\"2\", \"Smartphone\", \"Electronics\"],\n    [\"3\", \"Mouse\", \"Accessories\"],\n    [\"4\", \"Desk Chair\", \"Furniture\"],\n    [\"5\", \"Keyboard\", \"Accessories\"],\n    [\"6\", \"Monitor\", \"Electronics\"],\n    [\"7\", \"Book\", \"Books\"],\n    [\"8\", \"USB Cable\", \"Accessories\"],\n    [\"9\", \"Table\", \"Furniture\"],\n    [\"10\", \"Laptop Bag\", \"Accessories\"],\n    [\"11\", \"Notebook\", \"Stationery\"],\n    [\"12\", \"Pen\", \"Stationery\"],\n    [\"13\", \"Headphones\", \"Electronics\"],\n    [\"14\", \"Lamp\", \"Furniture\"],\n    [\"15\", \"Power Bank\", \"Electronics\"],\n    [\"16\", \"Backpack\", \"Accessories\"],\n    [\"17\", \"Tablet\", \"Electronics\"],\n    [\"18\", \"Coffee Mug\", \"Kitchenware\"],\n    [\"19\", \"Sofa\", \"Furniture\"],\n    [\"20\", \"Eraser\", \"Stationery\"]\n  ],\n  \"Orders\": [\n    [\"OrderID\", \"CustomerID\", \"OrderDate\"],\n    [\"1\", \"1\", \"2024-12-01\"],\n    [\"2\", \"2\", \"2024-12-02\"],\n    [\"3\", \"1\", \"2024-12-01\"],\n    [\"4\", \"3\", \"2024-12-03\"],\n    [\"5\", \"1\", \"2024-12-01\"],\n    [\"6\", \"2\", \"2024-12-02\"],\n    [\"7\", \"4\", \"2024-12-04\"],\n    [\"8\", \"1\", \"2024-12-01\"],\n    [\"9\", \"3\", \"2024-12-03\"],\n    [\"10\", \"2\", \"2024-12-02\"],\n    [\"11\", \"4\", \"2024-12-04\"],\n    [\"12\", \"1\", \"2024-12-01\"],\n    [\"13\", \"2\", \"2024-12-02\"],\n    [\"14\", \"3\", \"2024-12-03\"],\n    [\"15\", \"1\", \"2024-12-01\"],\n    [\"16\", \"4\", \"2024-12-04\"],\n    [\"17\", \"2\", \"2024-12-02\"],\n    [\"18\", \"1\", \"2024-12-01\"],\n    [\"19\", \"3\", \"2024-12-03\"],\n    [\"20\", \"4\", \"2024-12-04\"]\n  ],\n  \"OrderDetails\": [\n    [\"OrderID\", \"ProductID\", \"Quantity\", \"Price\"],\n    [\"1\", \"1\", \"1\", \"1200\"],\n    [\"2\", \"2\", \"2\", \"800\"],\n    [\"3\", \"3\", \"1\", \"25\"],\n    [\"4\", \"4\", \"1\", \"150\"],\n    [\"5\", \"5\", \"1\", \"45\"],\n    [\"6\", \"6\", \"1\", \"300\"],\n    [\"7\", \"7\", \"3\", \"30\"],\n    [\"8\", \"8\", \"2\", \"10\"],\n    [\"9\", \"9\", \"1\", \"200\"],\n    [\"10\", \"10\", \"1\", \"50\"],\n    [\"11\", \"11\", \"5\", \"15\"],\n    [\"12\", \"12\", \"10\", \"5\"],\n    [\"13\", \"13\", \"1\", \"100\"],\n    [\"14\", \"14\", \"2\", \"60\"],\n    [\"15\", \"15\", \"1\", \"40\"],\n    [\"16\", \"16\", \"1\", \"70\"],\n    [\"17\", \"17\", \"1\", \"500\"],\n    [\"18\", \"18\", \"2\", \"20\"],\n    [\"19\", \"19\", \"1\", \"800\"],\n    [\"20\", \"20\", \"10\", \"2\"]\n  ]\n}"
    return normalize_pb2.NormalizedTable(result=ret)

def serve():  
  server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
  normalize_pb2_grpc.add_NormalizeServiceServicer_to_server(NormalizeService(), server)
  server.add_insecure_port('[::]:50051')
  server.start()
  print("python grpc server running...")
  server.wait_for_termination()
  
  
  
if __name__ == '__main__':
  serve()