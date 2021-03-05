import Vue from 'vue';
import { colors } from 'vuetify/lib';
import Vuetify from 'vuetify/lib/framework';

Vue.use(Vuetify);

export default new Vuetify({
    theme: {
        dark: false,
        themes: {
            light: {
                primary: '#42A5F6',
                secondary: '#050B1F',
                accent: '#204165',
            },
            dark: {
                primary: '#50778D',
                secondary: '#0B1C3D',
                accent: '#204165',
            },
        }
    },
})
