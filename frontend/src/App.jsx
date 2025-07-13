import React, {useState} from 'react';
import {useEffect} from "react";
import {
    MenuList,
    MenuListItem,
    Separator,
    Frame,
    styleReset,
    Tabs,
    Tab,
    TabBody,
    TextInput,
    ScrollView,
    Button
} from 'react95';
import { createGlobalStyle, ThemeProvider } from 'styled-components';

import { MapContainer, TileLayer, useMap, Marker, Popup, Polyline} from 'react-leaflet'
import L from 'leaflet';
import original from 'react95/dist/themes/original';

import ms_sans_serif from './assets/fonts/fixedsys.woff2';
import ms_sans_serif_bold from 'react95/dist/fonts/ms_sans_serif_bold.woff2';
import RouteComponent from "./components/RouteComponent.jsx";
import BusComponent from "./components/BusComponent.jsx";
import DriverComponent from "./components/DriverComponent.jsx";
import BusStopComponent from "./components/BusStopComponent.jsx";




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




const App = () => {

    
    const [state, setState] = useState({
            activeTab: 0
      });
    const handleChange = (value, event) => {
            console.log({ value, event });
            setState({ activeTab: value });
          };
    const { activeTab } = state;
    return (
        <div>
            <GlobalStyles />
            <ThemeProvider theme={original}>
                <Frame
                    style={{ padding: '0.5rem', lineHeight: '1.5', width: 600 }}
                >
                    <Tabs value={"Маршруты"} onChange={handleChange}>
                        <Tab value={0}>Маршруты</Tab>
                        <Tab value={1}>Автобусы</Tab>
                        <Tab value={2}>Водители</Tab>
                        <Tab value={3}>Остановки</Tab>
                    </Tabs>
                    <TabBody>
                        {activeTab===0 && (
                            <RouteComponent />
                        )}
                        {activeTab===1 && (
                            <BusComponent />
                        )}
                        {activeTab===2 && (
                            <DriverComponent />
                        )}
                        {activeTab===3 && (
                            <BusStopComponent />
                        )}

                    </TabBody>

                </Frame>
            </ThemeProvider>
        </div>
    )
}


export default App;