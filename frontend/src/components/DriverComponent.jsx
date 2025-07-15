import {Button, MenuListItem, ScrollView, styleReset, TextInput} from "react95";
import React, {useEffect, useState} from "react";
import {createGlobalStyle} from "styled-components";
import ms_sans_serif from "../assets/fonts/fixedsys.woff2";
import ms_sans_serif_bold from 'react95/dist/fonts/ms_sans_serif_bold.woff2';
import original from 'react95/dist/themes/original';
import {GetAll, Add, DeleteById, GetById, UpdateById} from "../../wailsjs/go/routers/DriverRouter.js";
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

const DriverComponent = () => {
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
    const convertDateToISO = (date) => {
        if (!isValidDateFormat(date)) {
            return null; // Возвращаем null, если формат неверный
        }
        const [day, month, year] = date.split('.');
        return `${year}-${month.padStart(2, '0')}-${day.padStart(2, '0')}T00:00:00Z`;
    };
    const convertISOToDate = (isoDate) => {
        const date = new Date(isoDate);
        const day = date.getUTCDate().toString().padStart(2, '0');
        const month = (date.getUTCMonth() + 1).toString().padStart(2, '0');
        const year = date.getUTCFullYear().toString();
        return `${day}.${month}.${year}`;
    };
    const isValidDateFormat = (date) => {
        const regex = /^\d{2}\.\d{2}\.\d{4}$/;
        if (!regex.test(date)) return false;
        const [day, month, year] = date.split('.').map(Number);
        if (day < 1 || day > 31) return false;
        if (month < 1 || month > 12) return false;
        return !(year < 1900 || year > 9999);

    };
    const handleItemClick = (item) => {
        console.log("вадыла")

        GetById(item.ID).then(
            result => {
                const selectedData = JSON.parse(result);
                selectedData.BirthDate = convertISOToDate(selectedData.BirthDate)
                setSelectedItem(selectedData);
                console.log(selectedData)
            }
        ).catch(err => console.error("Ошибка при получении данных по ID:", err));
    };
    const handleInputChange = (field) => (e) => {
        setSelectedItem(prev => ({
            ...prev,
            [field]: e.target.value
        }));
    };
    const handleCloseAlert = () => {
        setAlertMessage(null);
    };
    const validateFields = (item) => {
        const fields = [
            { key: 'Surname', label: 'Фамилия' },
            { key: 'Name', label: 'Имя' },
            { key: 'Patronymic', label: 'Отчество' },
            { key: 'BirthDate', label: 'Дата рождения' },
            { key: 'PassportSeries', label: 'Серия и номер паспорта' },
            { key: 'Snils', label: 'Номер СНИЛС' },
            { key: 'LicenseSeries', label: 'Номер лицензии' }
        ];

        for (const field of fields) {
            if (!item[field.key] || item[field.key].trim() === '') {
                return { isValid: false, message: `Поле "${field.label}" не может быть пустым` };
            }
        }

        if (!isValidDateFormat(item.BirthDate)) {
            console.log(item)
            return { isValid: false, message: 'Дата должна быть в формате ДД.ММ.ГГГГ' };
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
            selectedItem.BirthDate = convertDateToISO(selectedItem.BirthDate)
            Add(JSON.stringify(selectedItem)).then(
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
            selectedItem.BirthDate = convertDateToISO(selectedItem.BirthDate)

            UpdateById(JSON.stringify(selectedItem)).then(
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
                        setAlertMessage(err);
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
            <ScrollView style={{width: '200px', height: '200px', marginRight: '20px'}}>
                {Array.isArray(items) && items.length > 0 ? (
                    items.map((item, index) => (
                        <MenuListItem key={index} onClick={() => handleItemClick(item)}>
                            {item.Surname}
                        </MenuListItem>
                    ))
                ) : (
                    <div>Нет водителей</div> // Отображение, если список пуст
                )}
            </ScrollView>
            <div>
                <div style={{ display: 'flex', flexDirection: 'row', alignItems: 'flex-start', marginBottom: '10px' }}>
                    <TextInput style={{width: '150px', marginRight: '20px'}}
                               value={selectedItem?.Name || ''}
                               onChange={handleInputChange('Name')}></TextInput>
                    <div style={{marginRight: '20px'}}>Имя</div>
                </div>
                <div style={{ display: 'flex', flexDirection: 'row', alignItems: 'flex-start', marginBottom: '10px' }}>
                    <TextInput style={{width: '150px', marginRight: '20px'}}
                               value={selectedItem?.Surname || ''}
                               onChange={handleInputChange('Surname')}></TextInput>
                    <div style={{marginRight: '20px'}}>Фамилия</div>
                </div>
                <div style={{ display: 'flex', flexDirection: 'row', alignItems: 'flex-start', marginBottom: '10px' }}>
                    <TextInput style={{width: '150px', marginRight: '20px'}}
                               value={selectedItem?.Patronymic || ''}
                        onChange={handleInputChange('Patronymic')}></TextInput>
                    <div style={{marginRight: '20px'}}>Отчество</div>
                </div>
                <div style={{ display: 'flex', flexDirection: 'row', alignItems: 'flex-start', marginBottom: '10px' }}>
                    <TextInput style={{width: '150px', marginRight: '20px'}}
                               value={selectedItem?.BirthDate || ''}
                               onChange={handleInputChange('BirthDate')}></TextInput>
                    <div style={{marginRight: '20px'}}>Дата рождения</div>
                </div>
                <div style={{ display: 'flex', flexDirection: 'row', alignItems: 'flex-start', marginBottom: '10px' }}>
                    <TextInput style={{width: '170px', marginRight: '20px'}}
                               value={selectedItem?.PassportSeries || ''}
                               onChange={handleInputChange('PassportSeries')}></TextInput>
                    <div style={{marginRight: '20px'}}>Серия и номер паспорта</div>
                </div>
                <div style={{ display: 'flex', flexDirection: 'row', alignItems: 'flex-start', marginBottom: '10px' }}>
                    <TextInput style={{width: '150px', marginRight: '20px'}}
                               value={selectedItem?.Snils || ''}
                               onChange={handleInputChange('Snils')}></TextInput>
                    <div style={{marginRight: '20px'}}>Номер СНИЛС</div>
                </div>
                <div style={{ display: 'flex', flexDirection: 'row', alignItems: 'flex-start', marginBottom: '10px' }}>
                    <TextInput style={{width: '150px', marginRight: '20px'}}
                               value={selectedItem?.LicenseSeries || ''}
                               onChange={handleInputChange('LicenseSeries')}></TextInput>
                    <div style={{marginRight: '20px'}}>Номер лицензии</div>
                </div>
                <div style={{ display: 'flex', flexDirection: 'row', alignItems: 'flex-start',  }}>
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

export default DriverComponent