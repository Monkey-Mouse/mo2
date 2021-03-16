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
        <v-expansion-panels v-model="expansion">
          <v-expansion-panel :key="0">
            <v-expansion-panel-header>
              <v-row>
                <v-col>
                  <h2>Light Theme</h2>
                </v-col>
              </v-row>
            </v-expansion-panel-header>
            <v-expansion-panel-content>
              <v-row><v-divider /></v-row>
              <v-row>
                <v-col>
                  <input-list
                    :inputProps="themeProps('light')"
                    :validator="validator"
                  />
                </v-col>
              </v-row>
            </v-expansion-panel-content>
          </v-expansion-panel>
          <v-expansion-panel :key="1">
            <v-expansion-panel-header>
              <v-row>
                <v-col>
                  <h2>Dark Theme</h2>
                </v-col>
              </v-row>
            </v-expansion-panel-header>
            <v-expansion-panel-content>
              <v-row><v-divider /></v-row>
              <v-row>
                <v-col>
                  <input-list
                    :inputProps="themeProps('dark')"
                    :validator="validator"
                  />
                </v-col>
              </v-row>
            </v-expansion-panel-content>
          </v-expansion-panel>
        </v-expansion-panels>
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
  themeColorChangeEx = new LazyExecutor();
  themeChangeEx = new LazyExecutor(() => {
    if (this.$refs["inputs"]) {
      (this.$refs["inputs"] as InputList).setModel({ darkMode: GetTheme() });
    }
  }, 1000);

  public get expansion() {
    return this.$vuetify.theme.dark ? 1 : 0;
  }

  public set expansion(v: number) {
    if (v === undefined) {
      return;
    }
    SetTheme(v === 1, this);
  }

  themeProps(theme: "light" | "dark") {
    const props: { [name: string]: InputProp } = {};
    for (const key in this.$vuetify.theme.themes[theme]) {
      props[key] = this.propFromThemeColor(key, "light");
    }
    return props;
  }

  propFromThemeColor(name: string, theme: "light" | "dark") {
    return {
      errorMsg: {},
      label: name.toUpperCase() + " COLOR",
      default: this.$vuetify.theme.themes[theme][name],
      icon: "mdi-account",
      col: 6,
      type: "color",
      onChange: (c: { hex: string }) => {
        this.themeColorChangeEx.Execute(
          () => (this.$vuetify.theme.themes[theme][name] = c.hex)
        );
      },
    };
  }
  @Watch("$vuetify", { immediate: true, deep: true })
  themeChange() {
    this.themeChangeEx.Execute();
  }
  beforeDestroy() {
    SetTheme(GetTheme(), this, this.$vuetify.theme.themes, this.user);
  }
}
</script>