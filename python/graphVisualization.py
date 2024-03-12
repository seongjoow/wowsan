import json
import networkx as nx
import matplotlib.pyplot as plt

folderPath = "data/graph"

# graph.json 파일에서 그래프 데이터 로드
filePath = folderPath + "/graph.json"
with open(filePath, 'r') as f:
    graph_data = json.load(f)

# 빈 그래프 객체 생성
G = nx.Graph()

# 노드와 엣지 추가
for node, neighbors in graph_data.items():
    G.add_node(node)
    for neighbor in neighbors:
        G.add_edge(node, neighbor)

# 그래프 시각화
plt.figure(figsize=(12, 12))
pos = nx.spring_layout(G)
nx.draw(G, pos, with_labels=True, node_color='skyblue', node_size=1000, edge_color='gray', linewidths=1, font_size=10)
# plt.axis('off')
plt.title("Graph Visualization")

# 이미지 파일 저장
image_path = folderPath + "/graph_visualization.png"
plt.savefig(image_path)

plt.show()