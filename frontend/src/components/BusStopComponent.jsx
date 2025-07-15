import {Button, MenuListItem, ScrollView, styleReset, TextInput} from "react95";
import React, {useEffect, useState} from "react";
import {createGlobalStyle} from "styled-components";
import ms_sans_serif from "../assets/fonts/fixedsys.woff2";
import ms_sans_serif_bold from 'react95/dist/fonts/ms_sans_serif_bold.woff2';
import { MapContainer, TileLayer, useMap, Marker, Popup, Polyline, useMapEvents} from 'react-leaflet'
import L from 'leaflet';
import original from 'react95/dist/themes/original';
import {Add, DeleteById, GetAll, GetById, UpdateById} from "../../wailsjs/go/routers/BusStopRouter.js";
import CustomAlert from "./CustomAlert.jsx";

const GlobalStyles = createGlobalStyle`
  ${styleReset}
  @font-face {
    font-family: 'ms_sans_serif';
    src: url('${ms_sans_serif}') format('woff2');
    font-weight: 400;
    font-style: normal
  }
  @font-face {
    font-family: 'ms_sans_serif';
    src: url('${ms_sans_serif_bold}') format('woff2');
    font-weight: bold;
    font-style: normal
  }
  body, input, select, textarea {
    font-family: 'ms_sans_serif';
  }
`;

const zoom = 15; // Уровень масштаба
let center = [53.23292, 44.87702]

const MapClickHandler = ({ onMapClick }) => {
    useMapEvents({
        click(e) {
            onMapClick(e.latlng); // Передаем координаты клика
        },
    });
    return null;
};

// Компонент для управления центром карты
const MapController = ({ center }) => {
    const map = useMap();
    useEffect(() => {
        if (center) {
            map.panTo(center);
        }
    }, [center, map]);
    return null;
};

const BusStopComponent = () => {
    const [markerPosition, setMarkerPosition] = useState(null);
    const [mapCenter, setMapCenter] = useState(center);

    const handleMapClick = (latlng) => {
        setMarkerPosition(latlng); // Сохраняем координаты маркера
        setSelectedItem(prev => ({
            ...prev,
            Lat: latlng.lat.toFixed(6), // Округляем до 6 знаков
            Long: latlng.lng.toFixed(6),
        }));
        console.log(selectedItem)
    };

    const [items, setItems] = useState([]);
    const [selectedItem, setSelectedItem] = useState(null);
    const [alertMessage, setAlertMessage] = useState(null);

    useEffect(() => {
        GetAll().then(
            result => {
                if (!result || result === "null") {
                    setItems([]);
                } else {
                    console.log(JSON.parse(result))
                    setItems(JSON.parse(result))
                }
            }
        )
    }, []);

    const handleItemClick = (item) => {
        console.log("остановка")

        GetById(item.ID).then(
            result => {
                const selectedData = JSON.parse(result);
                setSelectedItem(selectedData);
                console.log("GetById result:", selectedData);
                const lat = parseFloat(selectedData.Lat);
                const lng = parseFloat(selectedData.Long);
                const newCenter = [lat, lng];
                setMapCenter(newCenter);
                console.log("Попытка установки маркера:", { lat, lng, rawLat: selectedData.Lat, rawLong: selectedData.Long });

                const newPosition = { lat, lng };
                setMarkerPosition(newPosition);
                console.log(selectedData)
            }
        ).catch(err => console.error("Ошибка при получении данных по ID:", err));
    };

    const handleInputChange = (field) => (e) => {
        const value = e.target.value
        setSelectedItem(prev => ({
            ...prev,
            [field]: value
        }));
        if (field === 'Lat' || field === 'Long') {
            const lat = field === 'Lat' ? parseFloat(value) : (selectedItem?.Lat ? parseFloat(selectedItem.Lat) : 0);
            const lng = field === 'Long' ? parseFloat(value) : (selectedItem?.Long ? parseFloat(selectedItem.Long) : 0);
            if (!isNaN(lat) && !isNaN(lng)) {
                setMarkerPosition({ lat, lng });
                setMapCenter([lat, lng]);
            }
        }
    };

    const handleCloseAlert = () => {
        setAlertMessage(null);
    };

    const isValidCoordinateFormat = (value, isLatitude) => {
        if (!value || typeof value !== 'string') {
            console.log(`isValidCoordinateFormat: Неверный тип данных или пустая строка: ${value}`);
            return false;
        }

        // Проверяем, что строка содержит число (целое или с десятичной частью)
        const regex = /^-?\d+(\.\d+)?$/;
        if (!regex.test(value)) {
            console.log(`isValidCoordinateFormat: Неверный формат координаты: ${value}`);
            return false;
        }

        const num = parseFloat(value);
        if (isNaN(num)) {
            console.log(`isValidCoordinateFormat: Некорректное число: ${value}`);
            return false;
        }

        if (isLatitude) {
            if (num < -90 || num > 90) {
                console.log(`isValidCoordinateFormat: Широта вне диапазона [-90, 90]: ${value}`);
                return false;
            }
        } else {
            if (num < -180 || num > 180) {
                console.log(`isValidCoordinateFormat: Долгота вне диапазона [-180, 180]: ${value}`);
                return false;
            }
        }

        return true;
    };

    const validateFields = (item) => {
        const fields = [
            { key: 'Lat', label: 'Широта' },
            { key: 'Long', label: 'Долгота' },
            { key: 'Name', label: 'Название' },
        ];

        for (const field of fields) {
            const value = item[field.key] != null ? String(item[field.key]) : '';
            if (!value || value.trim() === '') {
                return { isValid: false, message: `Поле "${field.label}" не может быть пустым` };
            }
            if (field.isCoordinate && !isValidCoordinateFormat(item[field.key], field.isLatitude)) {
                return { isValid: false, message: `Поле "${field.label}" должно быть в формате десятичных градусов (${field.isLatitude ? 'широта: -90..90' : 'долгота: -180..180'})` };
            }
        }

        return { isValid: true, message: '' };
    };

    const handleCreate = () => {
        if (selectedItem) {
            const validation = validateFields(selectedItem);
            if (!validation.isValid) {
                setAlertMessage(validation.message);
                return;
            }
            selectedItem.ID = null
            const payload = {
                ...selectedItem,
                ID: null,
                Lat: parseFloat(selectedItem.Lat),
                Long: parseFloat(selectedItem.Long),
            };
            Add(JSON.stringify(payload)).then(
                result => {
                    if (JSON.parse(result).Error){
                        console.log(JSON.parse(result).Error)
                        setAlertMessage(JSON.parse(result).Error);
                    }

                    GetAll().then(
                        result => {
                            setItems(JSON.parse(result));
                            setSelectedItem(null);
                        }
                    ).catch(err => {
                        setAlertMessage(err);
                        console.error("Ошибка при обновлении списка:", err)
                    });
                }
            ).catch(err => {
                console.error("Ошибка при создании:", err)
                setAlertMessage(err);
            });
        } else {
            setAlertMessage("Нет выбранного элемента для создания");
            console.warn("Нет выбранного элемента для создания");
        }
    };

    const handleSave = () => {
        if (selectedItem) {
            if (!selectedItem.ID){
                setAlertMessage("Не выбран элемент");
                return;
            }
            const validation = validateFields(selectedItem);
            if (!validation.isValid) {
                setAlertMessage(validation.message);
                return;
            }

            const payload = {
                ...selectedItem,
                ID: selectedItem.ID,
                Lat: parseFloat(selectedItem.Lat),
                Long: parseFloat(selectedItem.Long),
            };
            UpdateById(JSON.stringify(payload)).then(
                result => {
                    if (JSON.parse(result).Error){
                        console.log(JSON.parse(result).Error)
                        setAlertMessage(JSON.parse(result).Error);
                    }

                    GetAll().then(
                        result => {
                            setItems(JSON.parse(result));
                            setSelectedItem(null);
                        }
                    ).catch(err => {
                        setAlertMessage(err);
                        console.error("Ошибка при обновлении списка:", err)
                    });
                }
            ).catch(err => {
                console.error("Ошибка при обновлении:", err)
                setAlertMessage(err);
            });
        } else {
            setAlertMessage("Нет выбранного элемента для обновления");
            console.warn("Нет выбранного элемента для обновления");
        }
    }

    const handleDelete = () => {
        if (selectedItem) {
            if (!selectedItem.ID){
                setAlertMessage("Нет выбранного элемента для удаления");
            }
            DeleteById(selectedItem.ID).then(
                result => {
                    GetAll().then(
                        result => {
                            setItems(JSON.parse(result));
                            setSelectedItem(null);
                        }
                    ).catch(err => {
                        setAlertMessage(Err);
                        console.error("Ошибка при обновлении списка:", err)
                    });
                }
            ).catch(err => {
                console.error("Ошибка при удалении:", err)
                setAlertMessage(err);
            });
        } else {
            setAlertMessage("Нет выбранного элемента для удаления");
            console.warn("Нет выбранного элемента для удаления");
        }
    }

    return (
        <div style={{ display: 'flex', flexDirection: 'row', alignItems: 'flex-start' }}>
            <GlobalStyles />
            <ScrollView style={{width: '200px', height: '200px', marginRight: '20px'}}>
                {Array.isArray(items) && items.length > 0 ? (
                    items.map((item, index) => (
                        <MenuListItem key={index} onClick={() => handleItemClick(item)}>
                            {item.Name}
                        </MenuListItem>
                    ))
                ) : (
                    <div>Нет остановок</div> // Отображение, если список пуст
                )}
            </ScrollView>
            <div>
                <div style={{ display: 'flex', flexDirection: 'row', alignItems: 'flex-start', marginBottom: '10px' }}>
                    <TextInput style={{width: '150px', marginRight: '20px'}}
                               value={selectedItem?.Lat || ''}
                               onChange={handleInputChange('Lat')}></TextInput>
                    <div style={{marginRight: '20px'}}>Широта</div>
                </div>
                <div style={{ display: 'flex', flexDirection: 'row', alignItems: 'flex-start', marginBottom: '10px' }}>
                    <TextInput style={{width: '150px', marginRight: '20px'}}
                               value={selectedItem?.Long || ''}
                               onChange={handleInputChange('Long')}></TextInput>
                    <div style={{marginRight: '20px'}}>Долгота</div>
                </div>
                <div style={{ display: 'flex', flexDirection: 'row', alignItems: 'flex-start', marginBottom: '10px' }}>
                    <TextInput style={{width: '150px', marginRight: '20px'}}
                               value={selectedItem?.Name || ''}
                               onChange={handleInputChange('Name')}></TextInput>
                    <div style={{marginRight: '20px'}}>Название</div>
                </div>
                <MapContainer
                    center={center}
                    zoom={zoom}
                    scrollWheelZoom={false}
                    style={{ height: '400px', width: '100%', marginTop: '20px', marginBottom: "20px" }}
                >
                    <TileLayer
                        attribution='© <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                    />
                    <MapClickHandler onMapClick={handleMapClick} />
                    <MapController center={mapCenter} />
                    {markerPosition && <Marker position={markerPosition}/>}
                </MapContainer>
                <div style={{ display: 'flex', flexDirection: 'row', alignItems: 'flex-start', marginRight: "20px"}}>
                    <Button style={{marginRight: '10px'}} onClick={handleSave}>Сохранить</Button>
                    <Button style={{marginRight: '10px'}} onClick={handleDelete}>Удалить</Button>
                    <Button onClick={handleCreate}>Создать</Button>
                </div>
            </div>
            {alertMessage && (
                <CustomAlert message={alertMessage} onClose={handleCloseAlert} />
            )}
        </div>
    )
}

export default BusStopComponent