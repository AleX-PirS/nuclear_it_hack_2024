# import matplotlib # type: ignore
# matplotlib.use('Agg')
# from flask import Flask, render_template, request, jsonify # type: ignore
# import geopandas as gpd # type: ignore
# import matplotlib.pyplot as plt # type: ignore
# from io import BytesIO
# import os

# app = Flask(__name__)

# @app.route('/')
# def index():
#     return render_template('index.html')

# @app.route('/upload', methods=['POST'])
# def upload():
#     # Получение данных из формы
#     file1 = request.files['file1']
#     file2 = request.files['file2']
#     accuracy = request.form['accuracy']
#     save_path = request.form['savePath']

#     # Чтение файлов GeoJSON из объектов BytesIO
#     gdf1 = gpd.read_file(BytesIO(file1.read()))
#     gdf2 = gpd.read_file(BytesIO(file2.read()))

#     # Настройка внешнего вида графика
#     fig, ax = plt.subplots(figsize=(20, 10))  # Устанавливаем размер графика

#     # Отображение данных из двух файлов на одном графике
#     gdf1.plot(ax=ax, color='blue', markersize=5, alpha=0.7, marker='o', label='GeoJSON Data 1', linestyle='-', linewidth=1)
#     gdf2.plot(ax=ax, color='red', markersize=5, alpha=0.7, marker='o', label='GeoJSON Data 2', linestyle='-', linewidth=1)
#     plt.legend()

#     # Сохранение графика в формате PNG
#     plot_image_path = 'static/img/plot.svg'
#     plt.savefig(plot_image_path, format='svg')


#     # Возвращаем страницу с графиком как изображением PNG
#     return render_template('plot.html', plot_image=plot_image_path)

# if __name__ == '__main__':
#     app.run(debug=True)


# Отображение файлов geojson с использованием plotly
import matplotlib
matplotlib.use('Agg')
from flask import Flask, render_template, request
import geopandas as gpd
import plotly.graph_objects as go
from io import BytesIO

app = Flask(__name__)

@app.route('/')
def index():
    return render_template('index.html')

@app.route('/upload', methods=['POST'])
def upload():
    # Получение данных из формы
    file1 = request.files['file1']
    file2 = request.files['file2']
    accuracy = request.form['accuracy']
    save_path = request.form['savePath']

    # Чтение файлов GeoJSON из объектов BytesIO
    gdf1 = gpd.read_file(BytesIO(file1.read()))
    gdf2 = gpd.read_file(BytesIO(file2.read()))

    # Обработка файлов GeoJSON и сохранение графика с использованием Plotly
    fig = process_geojson(gdf1, gdf2)

    # Преобразование карты в HTML для вставки в страницу
    plot_html = fig.to_html()

    # Возвращаем страницу с отображением карты
    return render_template('plotly_map.html', plot_html=plot_html)

def process_geojson(gdf1, gdf2):
    # Создание пустой карты
    fig = go.Figure()

    # Добавление MultiLineString объектов на карту
    for gdf, color in zip([gdf1, gdf2], ['blue', 'red']):
        for feature in gdf.iterfeatures():
            if feature['geometry']['type'] == 'MultiLineString':
                for line_coords in feature['geometry']['coordinates']:
                    lon, lat = zip(*line_coords)
                    fig.add_trace(go.Scattermapbox(
                        mode="lines",
                        lon=lon,
                        lat=lat,
                        line=dict(width=1, color=color),  # Используем разные цвета линий
                    ))

    # Настройка параметров карты
    fig.update_layout(
        mapbox=dict(
            style="open-street-map",
            center=dict(
                lon=gdf1.geometry.centroid.x.mean(),
                lat=gdf1.geometry.centroid.y.mean()
            ),
            zoom=20
        ),
        width=1720,  # Ширина карты в пикселях
        height=880,  # Высота карты в пикселях
    )

    return fig

if __name__ == '__main__':
    app.run(debug=True)
