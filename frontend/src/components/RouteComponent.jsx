import { Button, MenuListItem, ScrollView, styleReset, TextInput } from "react95";
import React, { useEffect, useRef, useState } from "react";
import { createGlobalStyle } from "styled-components";
import { MapContainer, TileLayer, Marker, Popup, useMap } from 'react-leaflet';
import {
    Add,
    DeleteById,
    GetAll,
    GetAllDriversById,
    GetAllBusStopsById,
    GetAllBusesById,
    AssignDriver,
    AssignBus,
    AssignBusStop,
    GetById,
    GetByNumber,
    UnassignDriver,
    UnassignBus,
    UnassignBusStop,
    UpdateById
} from "../../wailsjs/go/routers/RouteRouter.js";
import * as driverRouter from "../../wailsjs/go/routers/DriverRouter.js";
import * as busRouter from "../../wailsjs/go/routers/BusRouter.js";
import * as busStopRouter from "../../wailsjs/go/routers/BusStopRouter.js";
import L from 'leaflet';
import 'leaflet-routing-machine';
import original from 'react95/dist/themes/original';
import CustomAlert from "./CustomAlert.jsx";
import { createPortal } from "react-dom";
import ms_sans_serif from "../assets/fonts/fixedsys.woff2";
import ms_sans_serif_bold from 'react95/dist/fonts/ms_sans_serif_bold.woff2';
import { UpdateById as BusUpdateById } from "../../wailsjs/go/routers/BusRouter.js";

// Глобальные стили остаются без изменений
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

// Компонент Dropdown остаётся без изменений
const Dropdown = ({ items, onAdd, type, triggerRef }) => {
    const [isOpen, setIsOpen] = useState(false);
    const dropdownRef = useRef(null);

    useEffect(() => {
        const handleClickOutside = (event) => {
            if (dropdownRef.current && !dropdownRef.current.contains(event.target) && triggerRef.current && !triggerRef.current.contains(event.target)) {
                setIsOpen(false);
            }
        };
        document.addEventListener('mousedown', handleClickOutside);
        return () => {
            document.removeEventListener('mousedown', handleClickOutside);
        };
    }, [triggerRef]);

    const handleAdd = (item) => {
        onAdd(item);
        setIsOpen(false);
    };

    const renderDropdown = () => (
        isOpen && createPortal(
            <div
                ref={dropdownRef}
                style={{
                    position: 'absolute',
                    top: triggerRef.current ? triggerRef.current.getBoundingClientRect().bottom + window.scrollY : 0,
                    left: triggerRef.current ? triggerRef.current.getBoundingClientRect().left + window.scrollX : 0,
                    backgroundColor: '#c0c0c0',
                    border: '1px solid #000',
                    zIndex: 1000,
                    maxHeight: '150px',
                    overflowY: 'auto',
                    width: '150px',
                    boxShadow: '2px 2px 5px rgba(0,0,0,0.3)',
                }}
            >
                {items.length > 0 ? (
                    items.map((item, index) => (
                        <MenuListItem
                            key={index}
                            primary
                            size="sm"
                            onClick={() => handleAdd(item)}
                            style={{ width: '100%' }}
                        >
                            {type === 'driver' && item.Surname}
                            {type === 'busStop' && item.Name}
                            {type === 'bus' && item.RegisterNumber}
                        </MenuListItem>
                    ))
                ) : (
                    <MenuListItem primary size="sm" disabled>
                        Нет доступных элементов
                    </MenuListItem>
                )}
            </div>,
            document.body
        )
    );

    return (
        <>
            <MenuListItem
                primary
                size="sm"
                ref={triggerRef}
                onClick={() => setIsOpen(!isOpen)}
            >
                +
            </MenuListItem>
            {renderDropdown()}
        </>
    );
};

// Новый компонент для маршрутизации
const RoutingMachine = ({ busStops }) => {
    const map = useMap();
    const routingControlRef = useRef(null);

    useEffect(() => {
        if (!map || busStops.length < 2) return;

        // Очистка предыдущего маршрута
        if (routingControlRef.current) {
            map.removeControl(routingControlRef.current);
        }

        // Создание маршрута
        routingControlRef.current = L.Routing.control({
            waypoints: busStops.map(stop => L.latLng(stop.Lat, stop.Long)),
            lineOptions: {
                styles: [{ color: 'blue', weight: 4, opacity: 0.7 }],
            },
            router: L.Routing.osrmv1({
                serviceUrl: 'https://router.project-osrm.org/route/v1',
            }),
            show: false,
            addWaypoints: false,
            routeWhileDragging: false,
            fitSelectedRoutes: true,
            showAlternatives: false,
            showInstructions: false,
            createMarker: () => null,
        }).addTo(map);

        return () => {
            if (routingControlRef.current) {
                map.removeControl(routingControlRef.current);
            }
        };
    }, [map, busStops]);

    return null;
};

const RouteComponent = () => {
    const [items, setItems] = useState([]);
    const [selectedItem, setSelectedItem] = useState(null);
    const [alertMessage, setAlertMessage] = useState(null);
    const [drivers, setDrivers] = useState([]);
    const [busStops, setBusStops] = useState([]);
    const [buses, setBuses] = useState([]);
    const [availableDrivers, setAvailableDrivers] = useState([]);
    const [availableBusStops, setAvailableBusStops] = useState([]);
    const [availableBuses, setAvailableBuses] = useState([]);
    const mapRef = useRef(null);
    const driverTriggerRef = useRef(null);
    const busStopTriggerRef = useRef(null);
    const busTriggerRef = useRef(null);

    useEffect(() => {
        // Загрузка всех маршрутов
        GetAll().then(
            result => {

                if (!result || result === "null") {
                    setItems([]);
                    // setAlertMessage("Список маршрутов пуст");
                } else {
                    const parsed = JSON.parse(result);
                    if (parsed.length === 0) {
                        setItems([]);
                        setAlertMessage("Список маршрутов пуст");
                    } else {
                        setItems(parsed);
                    }
                }
            }
        ).catch(err => {
            // setAlertMessage("Ошибка при загрузке маршрутов: " + err);
            console.error("Ошибка при загрузке маршрутов:", err);
            setItems([]);
        });

        // Загрузка всех доступных водителей
        driverRouter.GetAll().then(result => {
            if (!result || result === "null") {
                setAvailableDrivers([]);
                // setAlertMessage("Ошибка загрузки водителей: данные отсутствуют");
            } else {
                const parsed = JSON.parse(result);
                if (parsed.Error) {
                    setAvailableDrivers([]);
                    // setAlertMessage("Ошибка загрузки водителей: " + parsed.Error);
                } else {
                    setAvailableDrivers(parsed);
                }
            }
        }).catch(err => {
            // setAlertMessage("Ошибка при загрузке водителей: " + err);
            console.error("Ошибка загрузки водителей:", err);
            setAvailableDrivers([]);
        });

        // Загрузка всех доступных остановок
        busStopRouter.GetAll().then(result => {
            if (!result || result === "null") {
                setAvailableBusStops([]);
                // setAlertMessage("Ошибка загрузки остановок: данные отсутствуют");
            } else {
                const parsed = JSON.parse(result);
                if (parsed.Error) {
                    setAvailableBusStops([]);
                    // setAlertMessage("Ошибка загрузки остановок: " + parsed.Error);
                } else {
                    setAvailableBusStops(parsed);
                }
            }
        }).catch(err => {
            // setAlertMessage("Ошибка при загрузке остановок: " + err);
            console.error("Ошибка при загрузке остановок:", err);
            setAvailableBusStops([]);
        });

        // Загрузка всех доступных автобусов
        busRouter.GetAll().then(result => {
            if (!result || result === "null") {
                setAvailableBuses([]);
                // setAlertMessage("Ошибка загрузки автобусов: данные отсутствуют");
            } else {
                const parsed = JSON.parse(result);
                if (parsed.Error) {
                    setAvailableBuses([]);
                    // setAlertMessage("Ошибка загрузки автобусов: " + parsed.Error);
                } else {
                    setAvailableBuses(parsed);
                }
            }
        }).catch(err => {
            // setAlertMessage("Ошибка при загрузке автобусов: " + err);
            console.error("Ошибка при загрузке автобусов:", err);
            setAvailableBuses([]);
        });
    }, []);

    const handleItemClick = (item) => {
        console.log("Выбран маршрут:", item);
        setSelectedItem(item);

        GetAllDriversById(item.ID).then(
            driverResult => {
                let driversData = [];
                if (JSON.parse(driverResult).Error) {
                    console.log(JSON.parse(driverResult));
                    setDrivers([]);
                    return;
                }
                driversData = JSON.parse(driverResult);
                console.log("Drivers:", driversData);
                setDrivers(driversData);
            }
        ).catch(err => {
            setAlertMessage("Ошибка при загрузке водителей: " + err);
            console.error("Ошибка при загрузке водителей:", err);
        });

        GetAllBusStopsById(item.ID).then(
            stopResult => {
                let stopsData = [];
                if (JSON.parse(stopResult).Error) {
                    setBusStops([]);
                    return;
                }
                stopsData = JSON.parse(stopResult).map(stop => ({
                    ...stop,
                    Lat: parseFloat(stop.Lat),
                    Long: parseFloat(stop.Long),
                }));
                console.log("Bus Stops:", stopsData);
                setBusStops(stopsData);
                if (stopsData.length > 0 && mapRef.current) {
                    mapRef.current.setView([stopsData[0].Lat, stopsData[0].Long], 15);
                }
            }
        ).catch(err => {
            setAlertMessage("Ошибка при загрузке остановок: " + err);
            console.error("Ошибка при загрузке остановок:", err);
        });

        GetAllBusesById(item.ID).then(
            busResult => {
                let busesData = [];
                if (JSON.parse(busResult).Error) {
                    setBuses([]);
                    return;
                }
                busesData = JSON.parse(busResult);
                console.log("Buses:", busesData);
                setBuses(busesData);
            }
        ).catch(err => {
            setAlertMessage("Ошибка при загрузке автобусов: " + err);
            console.error("Ошибка при загрузке автобусов:", err);
        });
    };

    const handleRemoveDriver = (driverToRemove) => {
        setDrivers((prevDrivers) => prevDrivers.filter((driver) => driver.ID !== driverToRemove.ID));
        setAvailableDrivers((prevAvailable) => [...prevAvailable, driverToRemove]);
    };

    const handleRemoveBusStop = (busStopToRemove) => {
        setBusStops((prevBusStops) => prevBusStops.filter((busStop) => busStop.ID !== busStopToRemove.ID));
        setAvailableBusStops((prevAvailable) => [...prevAvailable, busStopToRemove]);
    };

    const handleRemoveBus = (busToRemove) => {
        setBuses((prevBuses) => prevBuses.filter((bus) => bus.ID !== busToRemove.ID));
        setAvailableBuses((prevAvailable) => [...prevAvailable, busToRemove]);
    };

    const handleAddDriver = (driver) => {
        setDrivers((prevDrivers) => [...prevDrivers, driver]);
        setAvailableDrivers((prevAvailable) =>
            prevAvailable.filter((d) => d.ID !== driver.ID)
        );
        console.log("Добавлен водитель:", driver);
    };

    const handleAddBusStop = (busStop) => {
        setBusStops((prevBusStops) => [
            ...prevBusStops,
            { ...busStop, Lat: parseFloat(busStop.Lat), Long: parseFloat(busStop.Long) },
        ]);
        setAvailableBusStops((prevAvailable) =>
            prevAvailable.filter((s) => s.ID !== busStop.ID)
        );
        console.log("Добавлена остановка:", busStop);
    };

    const handleAddBus = (bus) => {
        setBuses((prevBuses) => [...prevBuses, bus]);
        setAvailableBuses((prevAvailable) =>
            prevAvailable.filter((b) => b.ID !== bus.ID)
        );
        console.log("Добавлен автобус:", bus);
    };

    const handleCloseAlert = () => {
        setAlertMessage(null);
    };

    const validate = () => {
        if (!Array.isArray(drivers) || drivers.length === 0) {
            return { isValid: false, message: 'Список водителей не может быть пустым' };
        }
        if (!Array.isArray(busStops) || busStops.length === 0) {
            return { isValid: false, message: 'Список остановок не может быть пустым' };
        }
        if (!Array.isArray(buses) || buses.length === 0) {
            return { isValid: false, message: 'Список автобусов не может быть пустым' };
        }
        if (!selectedItem || !selectedItem.Number || selectedItem.Number.trim() === '') {
            return { isValid: false, message: 'Номер маршрута не может быть пустым' };
        }
        return { isValid: true, message: '' };
    };

    const handleCreate = () => {
        if (selectedItem) {
            const valid = validate();
            if (!valid.isValid) {
                setAlertMessage(valid.message);
            } else {
                selectedItem.ID = null;
                Add(JSON.stringify(selectedItem)).then(
                    result => {
                        if (JSON.parse(result).Error) {
                            setAlertMessage(JSON.parse(result).Error);
                        }
                        GetAll().then(
                            result => {
                                setItems(JSON.parse(result));
                            }
                        ).catch(err => {
                            setAlertMessage(err);
                            console.error("Ошибка при обновлении списка:", err);
                        });
                    }
                ).catch(err => {
                    setAlertMessage(err);
                    console.error("Ошибка при создании:", err);
                });
                const number = selectedItem.Number;
                let id;
                GetByNumber(number).then(
                    result => {
                        id = JSON.parse(result).ID;
                        // Привязка водителей
                        drivers.forEach((element) => {
                            AssignDriver(id, element.ID).then(
                                result => {
                                    console.log("Водитель привязан:", result);
                                }
                            ).catch(err => {
                                setAlertMessage(err);
                                console.error("Ошибка при привязке водителя:", err);
                            });
                        });
                        // Привязка автобусов
                        buses.forEach((element) => {
                            AssignBus(id, element.ID).then(
                                result => {
                                    console.log("Автобус привязан:", result);
                                }
                            ).catch(err => {
                                setAlertMessage(err);
                                console.error("Ошибка при привязке автобуса:", err);
                            });
                        });
                        // Привязка остановок
                        busStops.forEach((element) => {
                            AssignBusStop(id, element.ID).then(
                                result => {
                                    console.log("Остановка привязана:", result);
                                }
                            ).catch(err => {
                                setAlertMessage(err);
                                console.error("Ошибка при привязке остановки:", err);
                            });
                        });
                    }
                ).catch(err => {
                    setAlertMessage(err);
                    console.error("Ошибка при получении маршрута:", err);
                });
                setSelectedItem(null);
            }
        } else {
            setAlertMessage("Нет выбранного элемента для создания");
            console.warn("Нет выбранного элемента для создания");
        }
    };

    const handleSave = () => {
        if (selectedItem) {
            const valid = validate();
            if (!valid.isValid) {
                setAlertMessage(valid.message);
                return;
            }

            const routeId = selectedItem.ID;
            console.log(JSON.stringify(selectedItem))
            // Сначала обновляем маршрут
            UpdateById(JSON.stringify(selectedItem))
                .then(result => {
                    const parsedResult = JSON.parse(result);
                    console.log(result)
                    if (parsedResult.Error) {
                        throw new Error(parsedResult.Error);
                    }
                    // Обещания для отвязки существующих сущностей только после успешного обновления
                    const unassignPromises = [
                        ...drivers.map(driver => UnassignDriver(routeId, driver.ID).catch(err => {
                            console.error(`Ошибка при отвязке водителя ${driver.ID}:`, err);
                            return Promise.reject(err);
                        })),
                        ...buses.map(bus => UnassignBus(routeId, bus.ID).catch(err => {
                            console.error(`Ошибка при отвязке автобуса ${bus.ID}:`, err);
                            return Promise.reject(err);
                        })),
                        ...busStops.map(busStop => UnassignBusStop(routeId, busStop.ID).catch(err => {
                            console.error(`Ошибка при отвязке остановки ${busStop.ID}:`, err);
                            return Promise.reject(err);
                        }))
                    ];
                    return Promise.all(unassignPromises);
                })
                .then(() => {
                    // Привязка новых сущностей
                    const assignPromises = [
                        ...drivers.map(driver => AssignDriver(routeId, driver.ID).catch(err => {
                            console.error(`Ошибка при привязке водителя ${driver.ID}:`, err);
                            return Promise.reject(err);
                        })),
                        ...buses.map(bus => AssignBus(routeId, bus.ID).catch(err => {
                            console.error(`Ошибка при привязке автобуса ${bus.ID}:`, err);
                            return Promise.reject(err);
                        })),
                        ...busStops.map(busStop => AssignBusStop(routeId, busStop.ID).catch(err => {
                            console.error(`Ошибка при привязке остановки ${busStop.ID}:`, err);
                            return Promise.reject(err);
                        }))
                    ];
                    return Promise.all(assignPromises);
                })
                .then(() => {
                    // Обновление списка и сброс состояний при успешном сохранении
                    return GetAll();
                })
                .then(result => {
                    setItems(JSON.parse(result));
                    setSelectedItem(null);
                    setDrivers([]);
                    setBusStops([]);
                    setBuses([]);
                })
                .catch(err => {
                    setAlertMessage(err.message || "Ошибка при сохранении маршрута");
                    console.error("Ошибка при сохранении маршрута:", err);
                });
        } else {
            setAlertMessage("Нет выбранного элемента для обновления");
            console.warn("Нет выбранного элемента для обновления");
        }
    };

    const handleDelete = () => {
        if (selectedItem.ID) {
            const routeId = selectedItem.ID;

            // Обещания для удаления связанных сущностей
            const unassignPromises = [
                ...drivers.map(driver => UnassignDriver(routeId, driver.ID).catch(err => {
                    console.error(`Ошибка при отвязке водителя ${driver.ID}:`, err);
                    return Promise.reject(err);
                })),
                ...buses.map(bus => UnassignBus(routeId, bus.ID).catch(err => {
                    console.error(`Ошибка при отвязке автобуса ${bus.ID}:`, err);
                    return Promise.reject(err);
                })),
                ...busStops.map(busStop => UnassignBusStop(routeId, busStop.ID).catch(err => {
                    console.error(`Ошибка при отвязке остановки ${busStop.ID}:`, err);
                    return Promise.reject(err);
                }))
            ];

            Promise.all(unassignPromises)
                .then(() => {
                    // Удаление маршрута после успешного отвязывания всех сущностей
                    return DeleteById(routeId);
                })
                .then(() => {
                    // Обновление списка и сброс состояний при успешном удалении
                    return GetAll();
                })
                .then(result => {
                    setItems(JSON.parse(result));
                    setSelectedItem(null);
                    setDrivers([]);
                    setBusStops([]);
                    setBuses([]);
                })
                .catch(err => {
                    setAlertMessage(err.message || "Ошибка при удалении маршрута");
                    console.error("Ошибка при удалении маршрута:", err);
                });
        } else {
            setAlertMessage("Нет выбранного элемента для удаления");
            console.warn("Нет выбранного элемента для удаления");
        }
    };

    const handleInputChange = (field) => (e) => {
        setSelectedItem(prev => ({
            ...prev,
            [field]: e.target.value
        }));
    };

    return (
        <div style={{ display: 'flex', flexDirection: 'row', alignItems: 'flex-start' }}>
            <GlobalStyles />
            <ScrollView style={{ width: '100px', height: '200px', marginRight: '20px' }}>
                {Array.isArray(items) && items.length > 0 ? (
                    items.map((item, index) => (
                        <MenuListItem key={index} onClick={() => handleItemClick(item)}>
                            {item.Name}
                        </MenuListItem>
                    ))
                ) : (
                    <div>Нет маршрутов</div> // Отображение, если список пуст
                )}
            </ScrollView>
            <div>
                <div style={{ display: 'flex', flexDirection: 'row', alignItems: 'flex-start', marginBottom: '10px' }}>
                    <TextInput
                        style={{ width: '150px', marginRight: '20px' }}
                        value={selectedItem ? selectedItem.Number : ''}
                        onChange={handleInputChange('Number')}
                    />
                    <div style={{ marginRight: '20px' }}>Номер маршрута</div>
                </div>
                <div style={{ display: 'flex', flexDirection: 'row', alignItems: 'flex-start', marginBottom: '10px' }}>
                    <ScrollView style={{ width: '150px', marginRight: '20px' }}>
                        {drivers.map((item, index) => (
                            <MenuListItem key={index} onClick={() => handleRemoveDriver(item)}>
                                {item.Surname}
                            </MenuListItem>
                        ))}
                    </ScrollView>
                    <div style={{ marginRight: '20px' }}>Водители на маршруте</div>
                    <Dropdown items={availableDrivers} onAdd={handleAddDriver} type="driver" triggerRef={driverTriggerRef} />
                </div>
                <div style={{ display: 'flex', flexDirection: 'row', alignItems: 'flex-start', marginBottom: '10px' }}>
                    <ScrollView style={{ width: '150px', marginRight: '20px' }}>
                        {busStops.map((item, index) => (
                            <MenuListItem key={index} onClick={() => handleRemoveBusStop(item)}>
                                {item.Name}
                            </MenuListItem>
                        ))}
                    </ScrollView>
                    <div style={{ marginRight: '20px' }}>Остановки</div>
                    <Dropdown items={availableBusStops} onAdd={handleAddBusStop} type="busStop" triggerRef={busStopTriggerRef} />
                </div>
                <div style={{ display: 'flex', flexDirection: 'row', alignItems: 'flex-start', marginBottom: '10px' }}>
                    <ScrollView style={{ width: '150px', marginRight: '20px' }}>
                        {buses.map((item, index) => (
                            <MenuListItem key={index} onClick={() => handleRemoveBus(item)}>
                                {item.RegisterNumber}
                            </MenuListItem>
                        ))}
                    </ScrollView>
                    <div style={{ marginRight: '20px' }}>Автобусы на маршруте</div>
                    <Dropdown items={availableBuses} onAdd={handleAddBus} type="bus" triggerRef={busTriggerRef} />
                </div>
                <div style={{ display: 'flex', flexDirection: 'row', alignItems: 'flex-start' }}>
                    <Button style={{ marginRight: '10px' }} onClick={handleSave}>Сохранить</Button>
                    <Button style={{ marginRight: '10px' }} onClick={handleDelete}>Удалить</Button>
                    <Button onClick={handleCreate}>Создать</Button>
                </div>
            </div>
            <MapContainer
                center={[53.23292, 44.87702]}
                zoom={15}
                scrollWheelZoom={false}
                style={{ height: '400px', width: '100%', marginTop: '20px' }}
                whenCreated={(map) => (mapRef.current = map)}
            >
                <TileLayer
                    attribution='© <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                    url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                />
                {busStops.map((stop, index) => (
                    <Marker key={index} position={[stop.Lat, stop.Long]}>
                        <Popup>{stop.Name}</Popup>
                    </Marker>
                ))}
                <RoutingMachine busStops={busStops} />
            </MapContainer>
            {alertMessage && (
                <CustomAlert message={alertMessage} onClose={handleCloseAlert} />
            )}
        </div>
    );
};

export default RouteComponent;