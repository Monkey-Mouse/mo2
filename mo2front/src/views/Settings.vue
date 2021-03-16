<template>
  <v-container>
    <v-row>
      <v-col cols="12" lg="8" offset-lg="2">
        <v-row>
          <v-col>
            <h1>Settings</h1>
          </v-col>
        </v-row>
        <v-row><v-divider /></v-row>
        <v-row>
          <v-col>
            <input-list
              ref="inputs"
              :inputProps="inputProps"
              :validator="validator"
            />
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <h2>Light Theme</h2>
          </v-col>
        </v-row>
        <v-row><v-divider /></v-row>
        <v-row>
          <v-col>
            <input-list :inputProps="inputPropslight" :validator="validator" />
          </v-col>
        </v-row>
        <v-row>
          <v-col>
            <h2>Dark Theme</h2>
          </v-col>
        </v-row>
        <v-row><v-divider /></v-row>
        <v-row>
          <v-col>
            <input-list :inputProps="inputPropsdark" :validator="validator" />
          </v-col>
        </v-row>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import { InputProp, User } from "@/models";
import { required } from "vuelidate/lib/validators";
import Vue from "vue";
import Component from "vue-class-component";
import InputList from "../components/InputList.vue";
import { GetTheme, LazyExecutor, SetTheme } from "@/utils";
import { Prop, Watch } from "vue-property-decorator";
@Component({
  components: { InputList },
})
export default class Settings extends Vue {
  @Prop()
  user: User;
  update = 0;
  validator = {};
  inputPropsdark: { [name: string]: InputProp } = {
    primary: {
      errorMsg: {},
      label: "Primary color",
      default: this.$vuetify.theme.themes.dark.primary,
      icon: "mdi-account",
      col: 6,
      type: "color",
      onChange: (c: { hex: string }) => {
        this.lightPrimaryChangeEx.Execute(
          () => (this.$vuetify.theme.themes.dark.primary = c.hex)
        );
      },
    },
    secondary: {
      errorMsg: {},
      label: "Secondary color",
      default: this.$vuetify.theme.themes.dark.secondary,
      icon: "mdi-account",
      col: 6,
      type: "color",
      onChange: (c: { hex: string }) => {
        this.lightSecondaryChangeEx.Execute(
          () => (this.$vuetify.theme.themes.dark.secondary = c.hex)
        );
      },
    },
  };
  inputPropslight: { [name: string]: InputProp } = {
    primary: {
      errorMsg: {},
      label: "Primary color",
      default: this.$vuetify.theme.themes.light.primary,
      icon: "mdi-account",
      col: 6,
      type: "color",
      onChange: (c: { hex: string }) => {
        this.lightPrimaryChangeEx.Execute(() => {
          this.$vuetify.theme.themes.light.primary = c.hex;
          console.log(this.$vuetify.theme.themes);
        });
      },
    },
    secondary: {
      errorMsg: {},
      label: "Secondary color",
      default: this.$vuetify.theme.themes.light.secondary,
      icon: "mdi-account",
      col: 6,
      type: "color",
      onChange: (c: { hex: string }) => {
        this.lightSecondaryChangeEx.Execute(
          () => (this.$vuetify.theme.themes.light.secondary = c.hex)
        );
      },
    },
  };
  inputProps: { [name: string]: InputProp } = {
    darkMode: {
      label: "prefer dark",
      default: GetTheme(),
      icon: "mdi-theme-light-dark",
      col: 6,
      type: "switch",
      errorMsg: {},
      onChange: () => SetTheme(!GetTheme(), this),
    },
  };
  lightPrimaryChangeEx = new LazyExecutor();
  lightSecondaryChangeEx = new LazyExecutor();
  themeChangeEx = new LazyExecutor(() => {
    if (this.$refs["inputs"]) {
      (this.$refs["inputs"] as InputList).setModel({ darkMode: GetTheme() });
    }
  }, 1000);
  @Watch("$vuetify", { immediate: true, deep: true })
  themeChange() {
    this.themeChangeEx.Execute();
  }
  beforeDestroy() {
    SetTheme(GetTheme(), this, this.$vuetify.theme.themes, this.user);
  }
}
</script>