<template>
  <v-avatar :size="size" color="brown">
    <v-img
      v-if="
        user.avatar !== null && user.avatar !== undefined && user.avatar !== ''
      "
      :src="user.avatar"
    ></v-img>
    <span v-else class="white--text headline">{{
      isUser ? initials : "A"
    }}</span>
  </v-avatar>
</template>

<script lang="ts">
import { User } from "@/models";
import { GetInitials, UserRole } from "@/utils";
import Vue from "vue";
import Component from "vue-class-component";
import { Prop } from "vue-property-decorator";
@Component({
  components: {},
})
export default class Avatar extends Vue {
  @Prop()
  user!: User;
  @Prop()
  size?: number;
  get isUser() {
    return this.user.roles && this.user.roles.indexOf(UserRole) >= 0;
  }
  get initials(): string {
    return GetInitials(this.user.name);
  }
}
</script>
