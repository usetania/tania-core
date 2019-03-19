<template lang="pug">
  .modal-template
    .modal-open
      .modal(style="display: block;")
        .modal-dialog
          .modal-content
            slot
              .modal-header
                span.h4.font-bold
                  translate This is the Modal Header
              .modal-body
                span.h4.font-bold
                  translate This is the Modal Body
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
  }
}
</script>

