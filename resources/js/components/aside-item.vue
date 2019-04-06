<template lang="pug">
  li(:class="active ? 'active': ''" @click="folded")
    slot
    slot(name="item")
</template>

<script>
export default {
  name: 'AsideItemComponent',
  data() {
    return {
      active: false,
    };
  },
  methods: {
    folded() {
      this.active = !this.active;
    },
  },
};
</script>
