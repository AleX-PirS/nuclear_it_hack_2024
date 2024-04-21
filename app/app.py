import os
from flask import Flask, render_template, request, jsonify, send_file  # type: ignore
import geopandas as gpd  # type: ignore
import requests  # type: ignore
from io import BytesIO  # type: ignore
import folium  # type: ignore

app = Flask(__name__)

# Определение папки для загрузки файлов
app.config['UPLOAD_FOLDER'] = os.path.join(os.path.dirname(__file__), 'uploads')

@app.route('/')
def index():
    return render_template('index.html')

@app.route('/upload', methods=['POST'])
def upload():

    GO_SERVICE_URL = 'http://localhost:8888/upload'

    # Проверка существования папки для сохранения файлов, если нет - создание
    if not os.path.exists(app.config['UPLOAD_FOLDER']):
        os.makedirs(app.config['UPLOAD_FOLDER'])
    
    # Получение файлов и значения точности из запроса
    file1 = request.files['file1']
    file2 = request.files['file2']
    accuracy = request.form['accuracy']
    
    # Сохранение файлов
    file1.save(os.path.join(app.config['UPLOAD_FOLDER'], file1.filename))
    file2.save(os.path.join(app.config['UPLOAD_FOLDER'], file2.filename))

    # Создание JSON с путями к сохраненным файлам и значением точности
    data = {
        'file1_path': os.path.join(app.config['UPLOAD_FOLDER'], file1.filename),
        'file2_path': os.path.join(app.config['UPLOAD_FOLDER'], file2.filename),
        'accuracy': accuracy
    }
    
    try:
        response = requests.post(GO_SERVICE_URL, json=data)  # Замена 'request' на 'requests'
        response.raise_for_status()  # Генерация ошибки при некорректных статусных кодах ответа
        
        # Обработка JSON-ответа от сервиса на Go
        file3 = response.json()['file3']  # Добавление квадратных скобок для доступа к ключу 'file3'
        file3_content = BytesIO(file3.encode())  # Преобразование содержимого ответа в объект BytesIO
        file3.save(os.path.join(app.config['UPLOAD_FOLDER'], file3.filename))
        # Здесь обрабатываем результат, например, создаем карту на основе данных из JSON

        # Возвращение ответа с данными, например, HTML с отображением карты
        return render_template('buttons.html', result=file3)
    except requests.exceptions.RequestException as e:  # Замена 'request' на 'requests'
        return jsonify({'error': str(e)}), 500  # Возвращение ответа с ошибкой


@app.route('/download_geojson', methods=['GET'])
def download_geojson():
    # Получение JSON-данных из запроса
    # Здесь необходимо реализовать логику для чтения данных GeoJSON из запроса
    # Предполагаем, что данные GeoJSON уже получены и сохранены в переменной 'file3'
    gdf1 = gpd.GeoDataFrame.from_features(file3['features']) # type: ignore
    temporary_geojson_path = "/path/to/temporary.geojson"
    gdf1.to_file(temporary_geojson_path, driver='GeoJSON')
    return send_file(temporary_geojson_path, as_attachment=True)
    # return "Download GEOJSON"  # Заглушка для примера


@app.route('/display_graph', methods=['GET', 'POST'])
def render_map():
    # Получение JSON-данных из запроса
    # result = response.json()

    # Чтение данных GeoJSON из JSON
    # Здесь необходимо реализовать логику для чтения данных GeoJSON из JSON
    # Предполагаем, что данные GeoJSON уже получены и сохранены в переменной 'file3'
    gdf1 = gpd.GeoDataFrame.from_features(file3['features']) # type: ignore

    # Создание карты с центром на центроиде геометрии
    m = folium.Map(location=[gdf1.geometry.centroid.y.mean(), gdf1.geometry.centroid.x.mean()], zoom_start=9)

    # Добавление GeoJSON на карту с всплывающими подсказками
    folium.GeoJson(gdf1, tooltip=folium.features.GeoJsonTooltip(fields=['id', 'road_name', 'start_m', 'finish_m'])).add_to(m)

    # Преобразование карты в HTML для вставки в страницу
    map_html = m._repr_html_()

    # Возвращение страницы с отображением карты
    return render_template('folium_map.html', map_html=map_html)

if __name__ == '__main__':
    app.run(debug=True, host="0.0.0.0", port="8080")
