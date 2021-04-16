import Vue from 'vue';
import Vuetify from 'vuetify/lib'

Vue.use(Vuetify);
export default new Vuetify({
    theme: {
        options: {
            customProperties: true
        },
        dark: false,
        themes: {
            light: {
                primary: '#607D8B',
                secondary: '#050B1F',
            },
            dark: {
                primary: '#50778D',
                secondary: '#0B1C3D',
                accent: '#474C50'
            },
        }
    },
})
