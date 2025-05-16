package ca.codeglacier.makishima.service;

import ca.codeglacier.makishima.entity.DiscordUser;
import ca.codeglacier.makishima.repository.DiscordUserRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class DiscordUserService {

    @Autowired
    private DiscordUserRepository userRepository;

    public void save(DiscordUser user) {
        userRepository.save(user);
    }

}
