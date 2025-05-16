package ca.codeglacier.makishima.entity;

import jakarta.persistence.Entity;
import jakarta.persistence.Id;

import java.util.Date;

@Entity
public class DiscordUser {

    @Id
    private Long id;
    private String accessToken;
    private String refreshToken;
    private Date tokenExpiry;

    public DiscordUser() {
    }

    public DiscordUser(Long id, String accessToken, String refreshToken, Date tokenExpiry) {
        this.id = id;
        this.accessToken = accessToken;
        this.refreshToken = refreshToken;
        this.tokenExpiry = tokenExpiry;
    }

}
