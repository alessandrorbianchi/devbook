$('#nova-mensagem').on('submit', enviarMensagem);

function enviarMensagem(evento) {
    evento.preventDefault();

    const elementoClicado = document.querySelector("#enviar-mensagem");
    const usuarioId = elementoClicado.getAttribute("data-usuario-id");
       
    $.ajax({
        url: `/usuarios/${usuarioId}/mensagens`,
        method: "POST",
        data: {
            mensagem: $('#mensagem').val(),
        }
    }).done(function() {
        window.location = `/mensagem-seguidor/${usuarioId}`;
    }).fail(function() {
        Swal.fire("Ops...", "Erro ao  enviar mensagem!", "error");
    })
} 