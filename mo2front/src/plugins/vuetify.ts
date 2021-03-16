import Vue from 'vue';
import { colors } from 'vuetify/lib';
import Vuetify from 'vuetify/lib/framework';

Vue.use(Vuetify);
export default new Vuetify({
    theme: {
        dark: false,
        themes: {
            light: {
                primary: '#607D8B',
                secondary: '#050B1F',
            },
            dark: {
                primary: '#50778D',
                secondary: '#0B1C3D',
            },
        }
    },
})
