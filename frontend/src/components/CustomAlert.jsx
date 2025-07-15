import {Button, MenuListItem, ScrollView, styleReset, TextInput, Window, WindowHeader, WindowContent} from "react95";
import React, {useEffect, useState} from "react";
import {createGlobalStyle} from "styled-components";
import ms_sans_serif from "../assets/fonts/fixedsys.woff2";
import ms_sans_serif_bold from 'react95/dist/fonts/ms_sans_serif_bold.woff2';
import original from 'react95/dist/themes/original';

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

const CustomAlert = ({ message, onClose }) => {
    return (
        <>
            {/* Полупрозрачный фон для блокировки интерфейса */}
            <div style={{
                position: 'fixed',
                top: 0,
                left: 0,
                right: 0,
                bottom: 0,
                backgroundColor: 'rgba(0, 0, 0, 0.5)',
                zIndex: 999,
            }} />
            {/* Окно алерта */}
            <div style={{
                position: 'fixed',
                top: '50%',
                left: '50%',
                transform: 'translate(-50%, -50%)',
                zIndex: 1000,
                width: '300px'
            }}>
                <Window>
                    <WindowHeader>Сообщение</WindowHeader>
                    <WindowContent>
                        <div style={{ marginBottom: '20px' }}>{message}</div>
                        <Button onClick={onClose} style={{ float: 'right' }}>OK</Button>
                    </WindowContent>
                </Window>
            </div>
        </>
    );
};

export default CustomAlert