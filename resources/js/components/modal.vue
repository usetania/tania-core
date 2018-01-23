<template lang="pug">
  .modal-open
    .modal.fade.in(style="display: block;")
      .modal-dialog
        .modal-content
          slot
            .modal-header
              span.h4.font-bold This is the Modal Header
            .modal-body
              span.h4.font-bold This is the Modal Body
    .modal-backdrop.fade.in
</template>

<script>
export default {
  name: 'AppModalComponent'
}
</script>

