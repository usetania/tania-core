import vClickOutside from 'v-click-outside'

<template lang="pug">
  .modal-template
    .modal-open
      .modal(style="display: block;")
        .modal-dialog
          .modal-content(v-click-outside="onClickOutside")
            slot
              .modal-header
                span.h4.font-bold This is the Modal Header
              .modal-body
                span.h4.font-bold This is the Modal Body
    .modal-backdrop.fade.in
</template>

<script>
export default {
  name: 'AppModalComponent',
  created: function () {
    var vm = this
    window.addEventListener('keyup', function(e) {
        if (e.keyCode == 27) {
          vm.$emit('close')
        }
    });
  },
  methods: {
    onClickOutside (event) {
      this.$emit('close')
    }
  }

}
</script>

